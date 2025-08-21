package server

// Do not touch this
// Routing config
// This file is a part of nadhi.dev/sarvar/fun


import "github.com/gofiber/fiber/v2"

var Route *fiber.App

func init() {
    Route = fiber.New()
}