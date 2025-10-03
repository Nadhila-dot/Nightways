package routes

import (
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"nadhi.dev/sarvar/fun/api-routes"
	"nadhi.dev/sarvar/fun/auth"
	"nadhi.dev/sarvar/fun/server"
)

// Import the extracted assets path
var assetsPath string
var ExtractedAssetsPath string

func SetAssetsPath(path string) {
    assetsPath = filepath.Join(path, "web/dist")
    ExtractedAssetsPath = path
}

func Register() {
	// Register all routes
	// Routes can be directly registered here or just put as functions
	z_inject()
	index()

	health()
	server.Route.Use("/api/v1", auth.CheckAuth)
	api.Index()
	api.KasmIndex()
	api.KasmCreateSession()
	api.AuthIndex()
	api.VelaIndex()
	api.SheetsIndex()
	api.RegisterWebsocketRoutes()
	api.Notebooks()
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

func index() {
	server.Route.Use(func(c *fiber.Ctx) error {
		path := c.Path()
		if !strings.HasPrefix(path, "/api/v1") && c.Method() == fiber.MethodGet {
			// Ignore SPA roots and all /vela/bucket routes
			if strings.HasPrefix(path, "/vela/bucket") || path == "/vela/" || path == "/index" || path == "/index.html" {
				return c.Next()
			}
			// Try to serve static file from extracted assets
			err := c.SendFile(filepath.Join(assetsPath, path))
			if err == nil {
				return nil
			}
			// If file not found, serve index.html for SPA
			return c.SendFile(filepath.Join(assetsPath, "index.html"))
		}
		return c.Next()
	})
}

func z_inject() {
    server.Route.Get("/z-inject/loader.js", func(c *fiber.Ctx) error {
        c.Set("Content-Type", "application/javascript")
        loaderPath := filepath.Join(ExtractedAssetsPath, "zp-inject/loader.js")
        return c.SendFile(loaderPath)
    })
}

func health() {
	server.Route.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "200",
			"message": "service is active",
		})
	})
}
