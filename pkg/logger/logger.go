// Package logger 提供写入本地文件的结构化日志，同时输出到标准输出。
// 日志文件按天滚动，存放在 logs/ 目录下，格式为 app-YYYY-MM-DD.log。
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const logDir = "logs"

var (
	mu      sync.Mutex
	current *os.File
	curDate string
	logger  *log.Logger
)

func init() {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		log.Fatalf("[logger] cannot create log dir: %v", err)
	}
	rotate()
}

// rotate 按当天日期切换日志文件（首次调用或跨天时触发）。
func rotate() {
	today := time.Now().Format("2006-01-02")
	if today == curDate && current != nil {
		return
	}
	if current != nil {
		_ = current.Close()
	}
	path := filepath.Join(logDir, fmt.Sprintf("app-%s.log", today))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("[logger] cannot open log file %s: %v", path, err)
	}
	current = f
	curDate = today
	// 同时写文件和标准输出
	mw := io.MultiWriter(os.Stdout, f)
	logger = log.New(mw, "", 0)
}

// write 是所有级别的底层写入函数。
func write(level, module, msg string, fields map[string]any) {
	mu.Lock()
	defer mu.Unlock()
	rotate()

	now := time.Now().Format("2006-01-02 15:04:05.000")
	line := fmt.Sprintf("[%s] [%s] [%s] %s", now, level, module, msg)
	if len(fields) > 0 {
		line += " |"
		for k, v := range fields {
			line += fmt.Sprintf(" %s=%v", k, v)
		}
	}
	logger.Println(line)
}

// Info 记录普通信息日志。
func Info(module, msg string, fields map[string]any) {
	write("INFO ", module, msg, fields)
}

// Warn 记录警告日志。
func Warn(module, msg string, fields map[string]any) {
	write("WARN ", module, msg, fields)
}

// Error 记录错误日志。
func Error(module, msg string, fields map[string]any) {
	write("ERROR", module, msg, fields)
}
