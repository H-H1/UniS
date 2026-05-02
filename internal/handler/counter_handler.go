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
	flushInterval = time.Minute // 最快每分钟落库一次
)

// CounterHandler 处理游客访问计数。
// 每次请求在内存中原子自增，每隔 flushInterval 才将累计值同步到数据库。
type CounterHandler struct {
	counterRepo repository.CounterRepository

	mu        sync.Mutex
	memCount  atomic.Int64 // 内存计数（所有请求累加，不受落库频率限制）
	lastFlush time.Time    // 上次落库时间
}

func NewCounterHandler(counterRepo repository.CounterRepository) *CounterHandler {
	h := &CounterHandler{
		counterRepo: counterRepo,
		lastFlush:   time.Now(),
	}

	// 启动时从数据库读取已有计数，作为内存计数的初始值
	if c, err := counterRepo.Get(counterURL); err == nil {
		h.memCount.Store(c.Count)
	}

	return h
}

// GetCounter godoc
// GET /counter
// 公开接口，无需登录。
// 每次请求内存计数 +1；每隔 1 分钟将当前内存值写入数据库。
func (h *CounterHandler) GetCounter(c *gin.Context) {
	// 1. 内存计数先加一
	current := h.memCount.Add(1)

	// 2. 判断是否需要落库（超过 1 分钟）
	h.mu.Lock()
	needFlush := time.Since(h.lastFlush) >= flushInterval*10
	if needFlush {
		h.lastFlush = time.Now()
	}
	h.mu.Unlock()

	if needFlush {
		if err := h.counterRepo.Save(counterURL, current); err != nil {
			logger.Error("counter_handler", "计数器落库失败", map[string]any{
				"url":       counterURL,
				"count":     current,
				"client_ip": c.ClientIP(),
				"error":     err.Error(),
			})
			// 落库失败不影响响应，继续返回内存计数
		} else {
			logger.Info("counter_handler", "计数器落库成功", map[string]any{
				"url":   counterURL,
				"count": current,
			})
		}
	}

	logger.Info("counter_handler", "游客访问计数", map[string]any{
		"url":       counterURL,
		"count":     current,
		"flushed":   needFlush,
		"client_ip": c.ClientIP(),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"url":   counterURL,
			"count": current,
		},
	})
}
