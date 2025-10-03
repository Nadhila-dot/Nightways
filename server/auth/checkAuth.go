package auth

import (
    "github.com/gofiber/fiber/v2"
    "strings"
)

func CheckAuth(c *fiber.Ctx) error {
    path := c.Path()
    // Allow /api/v1/auth*, /api/v1/info*, /api/v1/ws*, /api/v1/login*, and /api/v1/register* without auth
    if strings.HasPrefix(path, "/api/v1/auth") || strings.HasPrefix(path, "/api/v1/info") || strings.HasPrefix(path, "/api/v1/ws") || strings.HasPrefix(path, "/api/v1/login") || strings.HasPrefix(path, "/api/v1/register") {
        return c.Next()
    }

    authHeader := c.Get("Authorization")
    if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
        return c.Status(401).JSON(fiber.Map{"error": "missing or invalid authorization header"})
    }
    sessionID := authHeader[7:]
    valid, err := IsSessionValid(sessionID)
    if err != nil || !valid {
        return c.Status(401).JSON(fiber.Map{"error": "invalid session"})
    }
    return c.Next()
}

func GetUserInfo(c *fiber.Ctx) error {
    path := c.Path()
    // Allow /api/v1/auth*, /api/v1/info*, /api/v1/ws*, /api/v1/login*, and /api/v1/register* without auth (mirroring CheckAuth)
    if strings.HasPrefix(path, "/api/v1/auth") || strings.HasPrefix(path, "/api/v1/info") || strings.HasPrefix(path, "/api/v1/ws") || strings.HasPrefix(path, "/api/v1/login") || strings.HasPrefix(path, "/api/v1/register") {
        return c.Status(403).JSON(fiber.Map{"error": "authentication required for this endpoint"})
    }

    authHeader := c.Get("Authorization")
    if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
        return c.Status(401).JSON(fiber.Map{"error": "missing or invalid authorization header"})
    }
    sessionID := authHeader[7:]
    valid, err := IsSessionValid(sessionID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "server error validating session"})
    }
    if !valid {
        return c.Status(401).JSON(fiber.Map{"error": "invalid session"})
    }

    // Retrieve user info
    user, err := GetUserBySession(sessionID)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "user not found or session invalid"})
    }

    // Return user info as JSON (customize fields as needed)
    return c.JSON(user)
}

func getUsernameFromAuth(c *fiber.Ctx) (string, error) {
    authHeader := c.Get("Authorization")
    if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
        return "", fiber.ErrUnauthorized
    }
    sessionID := authHeader[7:]
    user, err := GetUserBySession(sessionID)
    if err != nil || user == nil {
        return "", fiber.ErrUnauthorized
    }
    return user.Username, nil
}