package api

import (
    "github.com/gofiber/fiber/v2"
    "nadhi.dev/sarvar/fun/server"
    "io/ioutil"
    "encoding/json"
   _ "time"
    ws "nadhi.dev/sarvar/fun/websocket"
)

func Index() error {
    server.Route.Get("/api/v1", func(c *fiber.Ctx) error {
        // Send broadcast notification when API index is hit
      //  ws.GetManager().Broadcast(ws.Info("info", "Someone accessed the API index"))
        
        return c.JSON(fiber.Map{
            "status": "200",
            "message": "service is active",
            "endpoints": []string{
                "/api/",
                "/api/kasm",
                "/api/v1/system",
                "/api/set (GET, POST)",
            },
        })
    })

    server.Route.Get("/api/v1/system", func(c *fiber.Ctx) error {
        // Send broadcast notification when system endpoint is hit
      //  ws.GetManager().Broadcast(ws.Info("system", "System information was requested"))
        
        return c.JSON(fiber.Map{
            "status": "200",
            "message": "service is active",
            "data": fiber.Map{
                "build":      "vela-0.1.0",
                "date":    "September 2025",
                "author":  "Nadhi",
            },
        })
    })

    // GET /api/set
    server.Route.Get("/api/v1/set", func(c *fiber.Ctx) error {
        filePath := "./set.json"
        dataBytes, err := ioutil.ReadFile(filePath)
        if err != nil {
            ws.GetManager().Broadcast(ws.Error("error", "Failed to read set.json", map[string]interface{}{}))
            return c.Status(500).JSON(fiber.Map{
                "error": "Failed to read set.json",
            })
        }
        
        var jsonData interface{}
        if err := json.Unmarshal(dataBytes, &jsonData); err != nil {
            ws.GetManager().Broadcast(ws.Error("error", "Invalid JSON in set.json", map[string]interface{}{}))
            return c.Status(500).JSON(fiber.Map{
                "error": "Invalid JSON in set.json",
            })
        }
        
        // Send broadcast notification when settings are retrieved
        //ws.GetManager().Broadcast(ws.Info("settings", "Settings were retrieved"))
        
        return c.JSON(jsonData)
    })

    // POST /api/set
    server.Route.Post("/api/v1/set", func(c *fiber.Ctx) error {
        auth := c.Get("Authorization")
        if len(auth) < 8 || auth[:7] != "Bearer " {
         //   ws.GetManager().Broadcast(ws.Warn("error", "Unauthorized attempt to update settings"))
            return c.Status(401).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }
        
        var newData interface{}
        if err := c.BodyParser(&newData); err != nil {
          //  ws.GetManager().Broadcast(ws.Warn("error", "Invalid JSON in settings update"))
            return c.Status(400).JSON(fiber.Map{
                "error": "Invalid JSON",
            })
        }
        
        filePath := "./set.json"
        dataBytes, err := json.MarshalIndent(newData, "", "  ")
        if err != nil {
          //  ws.GetManager().Broadcast(ws.Error("error", "Failed to encode JSON for settings", map[string]interface{}{}))
            return c.Status(500).JSON(fiber.Map{
                "error": "Failed to encode JSON",
            })
        }
        
        if err := ioutil.WriteFile(filePath, dataBytes, 0644); err != nil {
           // ws.GetManager().Broadcast(ws.Error("error", "Failed to write settings file", map[string]interface{}{}))
            return c.Status(500).JSON(fiber.Map{
                "status": 500,
                "error":  "Failed to write file",
            })
        }
        
        // Send broadcast notification when settings are updated
       // ws.GetManager().Broadcast(ws.Info("settings", "System settings were updated"))
        
        return c.JSON(fiber.Map{
            "status": 200,
            "message": "System environment updated successfully",
        })
    })

    return nil
}