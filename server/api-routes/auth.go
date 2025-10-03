package api

import (
	"github.com/gofiber/fiber/v2"
	"nadhi.dev/sarvar/fun/auth"
	"nadhi.dev/sarvar/fun/server"
	ws "nadhi.dev/sarvar/fun/websocket"
)

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Identifier string `json:"identifier"` // username or email
    Password   string `json:"password"`
}

func AuthIndex() error {
	// Register route
	server.Route.Post("/api/v1/register", func(c *fiber.Ctx) error {
		var body struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Rank     string `json:"rank"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}
		if body.Rank == "" {
			body.Rank = "user"
		}
		err := auth.Register(body.Username, body.Email, body.Password, body.Rank)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "registered"})
	})

	server.Route.Post("/api/v1/login", func(c *fiber.Ctx) error {
		var body struct {
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}
		sessionID, err := auth.Login(body.Identifier, body.Password)

		if err != nil {
			// Prob don't like this
			// It's not needed and is a waste of resources
			if err.Error() == "maximum sessions reached" {
				return c.Status(401).JSON(fiber.Map{"error": "maximum sessions reached"})
			}
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		if sessionID == "maximum sessions reached" {
			return c.Status(500).JSON(fiber.Map{"error": "maximum sessions reached"})
		}
		return c.Status(200).JSON(fiber.Map{"session": sessionID})
	})

	
	server.Route.Get("/api/v1/session/:id", func(c *fiber.Ctx) error {
    sessionID := c.Params("id")
    user, err := auth.GetUserBySession(sessionID)
    
    // Check for error first
    if err != nil {
        // Can't send to user directly since we don't know who they are
        // Just log the attempt
        c.Status(401)
        return c.JSON(fiber.Map{"error": "invalid session"})
    }
    
    // Now we know the user is valid, send them a notification
    ws.GetManager().Send(user.Username, ws.Info("info", "Session info requested"))
    
    return c.JSON(fiber.Map{
        "username": user.Username,
        "email":    user.Email,
        "rank":     user.Rank,
    })
})

	return nil
}
