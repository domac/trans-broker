package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (l *MyLogger) GinzapHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				l.zaplog.Error(e)
			}
		} else {

			status := c.Writer.Status()

			if status == http.StatusOK {
				l.zaplog.Info(path,
					zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.String("time", end.Format(API_LOG_TIME_FORMAT)),
					zap.Duration("latency", latency),
				)
			} else {
				l.zaplog.Warn(path,
					zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.String("time", end.Format(API_LOG_TIME_FORMAT)),
					zap.Duration("latency", latency),
				)
			}
		}
	}
}
