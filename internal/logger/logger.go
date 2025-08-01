package logger

import (
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func InitLogger() zerolog.Logger {
	graylogAddr := "localhost:12201"
	writer, err := gelf.NewWriter(graylogAddr)
	if err != nil {
		panic(err)
	}
	return zerolog.New(writer).With().Timestamp().Logger()
}

func LoggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = logger.With().Logger().WithContext(ctx)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
