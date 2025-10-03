package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"nadhi.dev/sarvar/fun/db"
	logg "nadhi.dev/sarvar/fun/logs"
	websocket "nadhi.dev/sarvar/fun/websocket"
)

var Route *fiber.App

func init() {
    Route = fiber.New()
    Route.Use(logger.New())
    logg.Info("Fiber instance created successfully")
    websocket.Init(log.Default())
    logg.Info("WebSocket manager initialized successfully")

    // Initialize databases
    // Sessions, Users and etc
    if err := db.InitSessionsDB(); err != nil {
        logg.Error("Failed to initialize sessions DB: ")
    }
    if err := db.InitUsersDB(); err != nil {
        logg.Error("Failed to initialize users DB: ")
    }
    if err := db.InitQueueDB(); err != nil {
    logg.Error("Failed to initialize queue DB: ")
    }
    if err := db.InitNotebooksDB(); err != nil {
    logg.Error("Failed to initialize notebooks DB: ")
    }
}