package routes

import (
    "nadhi.dev/sarvar/fun/server"
    "github.com/gofiber/fiber/v2"
)

func Register() {
	// Register all routes
	// Routes can be directly registered here or just put as functions
   	index()
	health()
}

/*
|--------------------------------------------------------------------------
| Routes for the App
|--------------------------------------------------------------------------
|
| Below this line, you can register your routes.
| You can also create separate functions for each route and call them
|
*/


func index(){
	server.Route.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("This is the index page!")
	})
}

func health(){
	server.Route.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "200",
			"message": "service is active",
		})
	})
}