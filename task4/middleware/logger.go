package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	logurs := logrus.New()
	logurs.SetFormatter(&logrus.JSONFormatter{})

	return gin.LoggerWithFormatter(func(p gin.LogFormatterParams) string {
		logurs.WithFields(logrus.Fields{
			"method":      p.Method,
			"timestamp":   p.TimeStamp.Format(time.DateTime),
			"status_code": p.StatusCode,
			"client_ip":   p.ClientIP,
			"path":        p.Path,
			"latency":     p.Latency,
		}).Info("HTTP Request ")
		return ""
	})
}

func ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error":     err,
					"path":      c.Request.URL.Path,
					"method":    c.Request.Method,
					"time":      c.Request.Header["time"],
					"client_ip": c.ClientIP(),
				}).Error("Panic recover")

				c.JSON(500, gin.H{
					"code":    500,
					"message": "Internal Server Error",
				})

				c.Abort()
			}
		}()
	}
}
