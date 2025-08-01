// @title 3x6 API
// @version 1.0
// @description API для игры "3 по 6"
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"github.com/beganov/gingonicserver/internal/api"
	"github.com/beganov/gingonicserver/internal/logger"
)

func main() {
	logger := logger.InitLogger()
	logger.Info().Msg("Logger initialized")
	router := api.SetupRouter(logger)
	if err := router.Run("localhost:8080"); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run server")
	}
}
