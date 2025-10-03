package api

import (
	"github.com/gofiber/fiber/v2"
	"nadhi.dev/sarvar/fun/auth"
	"nadhi.dev/sarvar/fun/kasm"
	"nadhi.dev/sarvar/fun/server"
)

func CheckAuth(c *fiber.Ctx) bool {
	sessionID := c.Query("sid")
	valid, err := auth.IsSessionValid(sessionID)
	if err != nil {
		return false
	}
	return valid
}


func KasmIndex() error {
	server.Route.Get("/api/kasm", func(c *fiber.Ctx) error {
		if !CheckAuth(c) {
			return c.JSON(fiber.Map{
				"status":  "401",
				"message": "unauthorized",
			})
		}
		return c.JSON(fiber.Map{
			"status": "200",
			"message": "All kasm endpoints will be listed here",
		})
	})
	return nil
}

func KasmCreateSession() {

	// Run function to create a new kasm session
	// Then return the kasm thing to login

	
	
	server.Route.Get("/api/kasm/create", func(c *fiber.Ctx) error {
		if !CheckAuth(c) {
			return c.JSON(fiber.Map{
				"status":  "401",
				"message": "unauthorized",
			})
		}
		kasmURL, err := kasm.GetLoginLink("0d2577de-8d25-44cd-94fc-2a05652c15a9")
		if err != nil {
			return c.JSON(fiber.Map{
					"status":  "500",
					"message": "internal server error",
					"error":   err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status": "200",
			"message": "Sucessfully created a new Kasm session",
			"url": kasmURL,
		})
	})
}
