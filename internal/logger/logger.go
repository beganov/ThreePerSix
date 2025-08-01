package logger

import (
	"os"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// InitLogger инициализирует zerolog.Logger, который пишет логи в Graylog через GELF-протокол
func InitLogger() zerolog.Logger {
	graylogAddr := "localhost:12201"           // Адрес сервера Graylog
	writer, err := gelf.NewWriter(graylogAddr) // Создаем GELF-объект для отправки логов
	if err != nil {
		zerolog.ConsoleWriter{Out: os.Stdout}.Write([]byte("Warning: failed to connect to Graylog, fallback to stdout\n"))
		return zerolog.Nop().With().Logger() // Возвращаем пустой логгер
	}
	return zerolog.New(writer).With().Timestamp().Logger()
}

func LoggerMiddleware(logger zerolog.Logger) gin.HandlerFunc { // LoggerMiddleware для добавления логгера в контекст запроса
	return func(c *gin.Context) {
		ctx := c.Request.Context()                    // Получаем контекст HTTP-запроса
		ctx = logger.With().Logger().WithContext(ctx) // встраиваем в него логгер
		c.Request = c.Request.WithContext(ctx)        // обновляем контекст запроса в объекте
		c.Next()                                      // переходим к следующему
	}
}
