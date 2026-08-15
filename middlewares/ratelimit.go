package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimiter 简单的内存滑动窗口限流器
type rateLimiter struct {
	mu      sync.Mutex
	limit   int           // 窗口内最大次数
	window  time.Duration // 时间窗口
	records map[string][]time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:   limit,
		window:  window,
		records: make(map[string][]time.Time),
	}
}

// allow 判断 key 是否允许通过，并记录本次访问
func (r *rateLimiter) allow(key string) bool {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := now.Add(-r.window)
	entries := r.records[key][:0]
	for _, t := range r.records[key] {
		if t.After(cutoff) {
			entries = append(entries, t)
		}
	}
	if len(entries) >= r.limit {
		r.records[key] = entries
		return false
	}
	r.records[key] = append(entries, now)

	// 防止内存无限增长: 记录过多时清理无近期访问的key
	if len(r.records) > 10000 {
		for k, v := range r.records {
			if len(v) == 0 {
				delete(r.records, k)
			}
		}
	}
	return true
}

// RateLimit 按客户端IP限流中间件
// 用法: RateLimit(10, 15*time.Minute) 表示每个IP每15分钟最多10次
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := newRateLimiter(limit, window)
	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{"msg": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}
		c.Next()
	}
}
