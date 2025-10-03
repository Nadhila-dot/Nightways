package api

import (
    "github.com/gofiber/fiber/v2"
    "nadhi.dev/sarvar/fun/server"
    "os"
    "path/filepath"
    "strings"
   _ "encoding/json"
    "io/ioutil"
)

// List all files in ./storage and send as JSON array (Fiber version)
func ListStorageFilesFiber(c *fiber.Ctx) error {
    files := []string{}
    root := "./storage"

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            relPath := strings.TrimPrefix(path, root+"/")
            files = append(files, relPath)
        }
        return nil
    })
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to list files"})
    }

    return c.JSON(files)
}

func ServeStorageFileFiber(c *fiber.Ctx) error {
    relPath := c.Params("*")
    filePath := filepath.Join("./storage", relPath)

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "File not found"})
    }

    raw := c.Query("raw") == "true"

    switch ext := strings.ToLower(filepath.Ext(filePath)); ext {
    case ".png":
        c.Type("png")
    case ".jpg", ".jpeg":
        c.Type("jpeg")
    case ".pdf":
        c.Type("pdf")
    case ".txt", ".md", ".json", ".csv":
        // Force text display in browser like GitHub Raw
        c.Set("Content-Type", "text/plain; charset=utf-8")
        c.Set("Content-Disposition", "inline")
        c.Set("X-Content-Type-Options", "nosniff")
        
        if raw {
            return c.Send(data)
        }
        return c.SendString(string(data))
    default:
		c.Set("Content-Type", "text/plain; charset=utf-8")
        c.Set("Content-Disposition", "inline")
        c.Set("X-Content-Type-Options", "nosniff")
        
    }

    return c.Send(data)
}

// Register Vela storage routes
func VelaIndex() error {
    // List all files in ./storage as JSON array
    server.Route.Get("/vela/list", ListStorageFilesFiber)

    // Serve a file from ./storage/bucket/whatever
    server.Route.Get("/vela/bucket/*", ServeStorageFileFiber)

    return nil
}