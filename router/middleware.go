package router

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	tlog "github.com/yohobala/taurus_go/tlog"
)

// 请求信息日志中间件
//
// 会记录 请求时间、请求的ip、请求的url、请求的方法、状态码、响应时间
func InfoLoggerMiddleware(log *tlog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		t := time.Now()
		ip := getRealIP(c)
		// 在 gin 上下文中定义变量
		c.Set("rquestTime", t)
		c.Set("ip", ip)
		// 请求前
		c.Next()

		// 请求后
		latency := time.Since(t)

		logger := tlog.Get("api")
		log, exists := c.Get("log")
		if !exists {
			logger.Debug("",
				tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")), // 请求时间
				tlog.String("ip", ip),                    // 请求的ip
				tlog.String("url", c.Request.URL.Path),   // 请求的url
				tlog.String("method", c.Request.Method),  // 请求的方法
				tlog.Int("code", c.Writer.Status()),      // 状态码
				tlog.String("latency", latency.String()), // 响应时间
			) // Structured context as strongly typed Field values.
		} else {
			logger.Debug("",
				tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")), // 请求时间
				tlog.String("ip", ip),                    // 请求的ip
				tlog.String("url", c.Request.URL.Path),   // 请求的url
				tlog.String("method", c.Request.Method),  // 请求的方法
				tlog.Int("code", c.Writer.Status()),      // 状态码
				tlog.String("latency", latency.String()), // 响应时间
				tlog.String("log", log.(string)),         // 响应时间
			) // Structured context as strongly typed Field values.
		}

	}
}

func getRealIP(c *gin.Context) string {
	// Try to get IP from X-Real-IP header
	clientIP := c.Request.Header.Get("X-Real-IP")
	if clientIP != "" {
		return clientIP
	}

	// Fallback to X-Forwarded-For header
	clientIP = c.Request.Header.Get("X-Forwarded-For")
	if clientIP != "" {
		return strings.Split(clientIP, ",")[0]
	}

	// Fallback to request remote address
	return c.ClientIP()
}
