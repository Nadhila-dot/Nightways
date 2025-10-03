package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"nadhi.dev/sarvar/fun/auth"
	store "nadhi.dev/sarvar/fun/database"
	"nadhi.dev/sarvar/fun/db"
	notebook "nadhi.dev/sarvar/fun/notebooks"
	"nadhi.dev/sarvar/fun/server"
)

func UpdateNotebook(username string, id int, name, description string, optional notebook.Optional) (*store.Notebook, error) {
    nb, err := store.GetNotebook(db.NotebooksDB, username, id)
    if err != nil {
        return nil, err
    }

    nb.Name = name
    nb.Description = description
    nb.Optional = store.Optional{
        Tags:  optional.Tags,
        Color: optional.Color,
    }

    err = store.UpdateNotebook(db.NotebooksDB, username, *nb)
    if err != nil {
        return nil, err
    }

    return nb, nil
}

// Helper to get username from session
func getUsernameFromAuth(c *fiber.Ctx) (string, error) {
    authHeader := c.Get("Authorization")
    if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
        return "", fiber.ErrUnauthorized
    }
    sessionID := authHeader[7:]
    user, err := auth.GetUserBySession(sessionID)
    if err != nil || user == nil {
        return "", fiber.ErrUnauthorized
    }
    return user.Username, nil
}

func Notebooks() error {
    server.Route.Get("/api/v1/notebooks", func(c *fiber.Ctx) error {
        username, err := getUsernameFromAuth(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        notebooks, err := notebook.GetAllNotebooks(username)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to get notebooks"})
        }
        return c.JSON(notebooks)
    })

    server.Route.Post("/api/v1/notebooks", func(c *fiber.Ctx) error {
        username, err := getUsernameFromAuth(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        var body struct {
            Name        string   `json:"name"`
            Description string   `json:"description"`
            Tags        []string `json:"tags"`
            Color       string   `json:"color"`
        }
        if err := c.BodyParser(&body); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
        }
        nb, err := notebook.CreateNotebook(username, body.Name, body.Description, notebook.Optional{
            Tags:  body.Tags,
            Color: body.Color,
        })
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to create notebook"})
        }
        return c.JSON(nb)
    })

    server.Route.Get("/api/v1/notebooks/:id", func(c *fiber.Ctx) error {
        username, err := getUsernameFromAuth(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        id, err := c.ParamsInt("id")
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid notebook id"})
        }
        nb, err := notebook.GetNotebook(username, id)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "notebook not found"})
        }
        return c.JSON(nb)
    })

    server.Route.Get("/api/v1/notebooks/:id/items", func(c *fiber.Ctx) error {
        username, err := getUsernameFromAuth(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        id, err := c.ParamsInt("id")
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid notebook id"})
        }
        items, err := notebook.GetItemsInNotebook(username, id)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "notebook not found"})
        }
        return c.JSON(items)
    })

	server.Route.Delete("/api/v1/notebooks/:id", func(c *fiber.Ctx) error {
		username, err := getUsernameFromAuth(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid notebook id"})
		}
		err = notebook.DeleteNotebook(username, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to delete notebook"})
		}
		return c.JSON(fiber.Map{"status": "deleted"})
	})

	server.Route.Put("/api/v1/notebooks/:id", func(c *fiber.Ctx) error {
    username, err := getUsernameFromAuth(c)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
    }
    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid notebook id"})
    }
    var body struct {
        Name        string   `json:"name"`
        Description string   `json:"description"`
        Tags        []string `json:"tags"`
        Color       string   `json:"color"`
    }
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
    }

    // Remove unused nb variable and call UpdateNotebook directly
    updatedNb, err := UpdateNotebook(username, id, body.Name, body.Description, notebook.Optional{
		Tags:  body.Tags,
		Color: body.Color,
	})
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to update notebook"})
    }

    return c.JSON(updatedNb)
})

	server.Route.Delete("/api/v1/notebooks/:id/items/:itemName", func(c *fiber.Ctx) error {
		username, err := getUsernameFromAuth(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid notebook id"})
		}
		itemName := c.Params("itemName")
		
		err = notebook.DeleteItemFromNotebook(username, id, itemName)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to remove sheet"})
		}
		return c.JSON(fiber.Map{"status": "removed"})
	})


	server.Route.Post("/api/v1/notebooks/:id/items", func(c *fiber.Ctx) error {
		username, err := getUsernameFromAuth(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid notebook id"})
		}
		var body struct {
			SheetName string `json:"sheetName"`
			Url       string `json:"url"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}
		err = notebook.CreateItemToNotebook(username, id, body.SheetName, body.Url)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to add sheet"})
		}
		return c.JSON(fiber.Map{"status": "added"})
	})

    return nil
}