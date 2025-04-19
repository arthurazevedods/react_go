package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/teste", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Teste")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))

}
