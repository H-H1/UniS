package handler

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"uniS/internal/repository"
	"uniS/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	counterURL    = "/counter"
	flushInterval = 10 * time.Minute // 每 10 分钟落库一次
)

// CounterHandler 处理游客访问计数。
//
// 流程：
//  1. 每次 GET /counter 请求，内存原子计数 +1，立即返回当前总量。
//  2. 每隔 10 分钟，将这段时间的增量作为一条新记录插入数据库。
//  3. 响应中的 count 是内存总量（历史落库总和 + 本轮未落库增量）。
type CounterHandler struct {
	counterRepo repository.CounterRepository

	mu         sync.Mutex
	increment  atomic.Int64 // 本轮未落库的增量
	totalCount atomic.Int64 // 内存中的总量（历史 + 本轮）
	lastFlush  time.Time    // 上次落库时间
}

func NewCounterHandler(counterRepo repository.CounterRepository) *CounterHandler {
	h := &CounterHandler{
		counterRepo: counterRepo,
		lastFlush:   time.Now(),
	}

	// 启动时从数据库读取历史总量，作为内存总量的基准
	if total, err := counterRepo.Total(counterURL); err == nil {
		h.totalCount.Store(total)
	} else {
		logger.Error("counter_handler", "启动时读取历史总量失败", map[string]any{
			"url":   counterURL,
			"error": err.Error(),
		})
	}

	return h
}

// GetCounter godoc
// GET /counter
// 公开接口，无需登录。
// 每次请求内存计数 +1；每满 10 分钟将本轮增量插入数据库一条新记录。
func (h *CounterHandler) GetCounter(c *gin.Context) {
	// 1. 内存总量和本轮增量同时 +1
	total := h.totalCount.Add(1)
	h.increment.Add(1)

	// 2. 判断是否到了落库时间
	h.mu.Lock()
	needFlush := time.Since(h.lastFlush) >= flushInterval
	var delta int64
	if needFlush {
		delta = h.increment.Swap(0) // 取出增量并重置为 0
		h.lastFlush = time.Now()
	}
	h.mu.Unlock()

	// 3. 落库（在锁外执行，不阻塞后续请求）
	if needFlush && delta > 0 {
		if err := h.counterRepo.Insert(counterURL, delta); err != nil {
			logger.Error("counter_handler", "计数器落库失败", map[string]any{
				"url":       counterURL,
				"delta":     delta,
				"total":     total,
				"client_ip": c.ClientIP(),
				"error":     err.Error(),
			})
			// 落库失败：把增量加回去，等下次再试
			h.increment.Add(delta)
		} else {
			logger.Info("counter_handler", "计数器落库成功", map[string]any{
				"url":   counterURL,
				"delta": delta, // 本轮增量
				"total": total, // 内存总量
			})
		}
	}

	logger.Info("counter_handler", "游客访问计数", map[string]any{
		"url":       counterURL,
		"total":     total,
		"flushed":   needFlush,
		"client_ip": c.ClientIP(),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"url":   counterURL,
			"count": total,
		},
	})
}
