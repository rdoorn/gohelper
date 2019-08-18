package ginhelper

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rdoorn/gohelper/logging"
)

func Logger(logging logging.SimpleLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		end := time.Now()
		if raw != "" {
			path = path + "?" + raw
		}

		logging.Infof(
			//path,
			c.Request.Host,
			"client", c.ClientIP(),
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", end.Sub(start),
			"size", c.Writer.Size(),
			"error", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)

	}
}
