package api

import (
	_ "github.com/beganov/gingonicserver/internal/api/docs"
	"github.com/beganov/gingonicserver/internal/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RouteRegister(router *gin.Engine) { // RouteRegister регистрирует все маршруты API в роутере Gin
	server := NewServer()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Маршрут для сваггера

	// Роуты для управления комнатами игры

	// Создать новую комнату
	// Выход: JSON с полями roomId и playerId (хост)
	router.POST("/rooms/", server.createRoom)

	// Получить данные комнаты по ID
	// :id - id комнаты
	// Выход: JSON с объектом room.Room
	router.GET("/rooms/:id/", server.getRoom)

	// Удалить комнату
	// :id - id комнаты
	router.DELETE("/rooms/:id/", server.deleteRoom)

	// Обновить maxPlayers комнаты
	// :id - id комнаты
	// Вход: JSON с объектом room.RoomUpdate
	router.PATCH("/rooms/:id/", server.patchRoom)

	// Присоединиться к комнате
	// :id - id комнаты
	// Выход: JSON с playerId
	router.POST("/rooms/:id/join", server.joinRoom)

	// Выйти из комнаты
	// :id - id комнаты
	// Вход: JSON с playerId
	router.DELETE("/rooms/:id/leave", server.leaveRoom)

	// Запустить игру
	// :id - id комнаты
	// Выход: JSON с объектом room.Room
	router.POST("/rooms/:id/start", server.start)

	// Сделать ход в игре
	// :id - id комнаты
	// Вход: JSON с полями playerId и move
	router.POST("/rooms/:id/move", server.move)
}

func SetupRouter(zllogger zerolog.Logger) *gin.Engine {
	router := gin.Default() // Создаем дефолтный роутер

	router.Use(cors.New(cors.Config{ // Настройка CORS для разрешения запросов с фронтенда на localhost:3000
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.Use(logger.LoggerMiddleware(zllogger)) // Подключаем middleware логгирования с переданным zerolog.Logger

	RouteRegister(router) // Регистрирация маршрутов API

	return router
}
