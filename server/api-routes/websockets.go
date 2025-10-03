package api

import (
	"crypto/md5"
	"fmt"
	_ "time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"nadhi.dev/sarvar/fun/auth"
	"nadhi.dev/sarvar/fun/server"

	sheet "nadhi.dev/sarvar/fun/sheets"
	ws "nadhi.dev/sarvar/fun/websocket"
)

func RegisterWebsocketRoutes() {
	// Middleware to check if connection is websocket
	server.Route.Use("/api/v1/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Websocket connection endpoint
	server.Route.Get("/api/v1/ws/notifications", websocket.New(ws.GetManager().ConnectHandler))

	// Websocket status endpoint (HTTP)
	server.Route.Get("/api/v1/ws/status", func(c *fiber.Ctx) error {
		return c.JSON(ws.GetManager().Status())
	})

	server.Route.Get("/api/v1/ws/job/:jobid", websocket.New(func(c *websocket.Conn) {
		jobID := c.Params("jobid")
		sessionID := c.Query("session")

		// Validate session
		isValid, err := auth.IsSessionValid(sessionID)
		if err != nil || !isValid {
			c.WriteJSON(ws.Error(
				"Invalid session",
				"Authentication failed",
				map[string]interface{}{},
			))
			c.Close()
			return
		}

		// Check if GlobalSheetGenerator is initialized
		if sheet.GlobalSheetGenerator == nil || sheet.GlobalSheetGenerator.Queue == nil {
			c.WriteJSON(ws.Error(
				"Server error",
				"Sheet generator not initialized",
				map[string]interface{}{},
			))
			c.Close()
			return
		}

		// Create a map to track last sent status per job
		lastSent := make(map[string]string)

		// Register a listener for this job
		sheet.GlobalSheetGenerator.Queue.RegisterJobListener(jobID, func(update sheet.StatusUpdate) {
			// Build a string to hash only the relevant fields
			hashInput := fmt.Sprintf("%s|%v|%v", update.Status, update.Result, update.Data)
			hash := fmt.Sprintf("%x", md5.Sum([]byte(hashInput)))

			// If we've already sent this status, skip it
			if lastSent[jobID] == hash {
				return
			}
			lastSent[jobID] = hash

			// Create a message with the update data
			msg := ws.Push(update.Status, map[string]interface{}{
				"message": fmt.Sprintf("Job %s status: %s", update.ID, update.Status),
			})
			msg["jobId"] = update.ID
			if update.Result != nil {
				msg["result"] = update.Result
			}
			if update.Data != nil {
				msg["data"] = update.Data
			}

			// Send to the client
			_ = c.WriteJSON(msg)
		})

		// Optionally send initial status
		job, exists := sheet.GlobalSheetGenerator.Queue.GetJobStatus(jobID)
		if exists {
			message := fmt.Sprintf("Initial status for job %s: %s", jobID, job.Status)
			msg := ws.Start(message, map[string]interface{}{})
			msg2 := ws.Info(job.Status, message)
			msg["jobId"] = jobID
			c.WriteJSON(msg)
			c.WriteJSON(msg2)
		}

		// Keep connection open until it's closed by the client
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				// Connection is closed
				break
			}
		}
	}))
}
