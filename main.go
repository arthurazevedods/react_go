package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	todos := []Todo{}

	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/hello", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "hello postman"})
	})

	app.Get("/teste", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Teste")
	})

	app.Post("/api/todos/", func(c fiber.Ctx) error {
		todo := &Todo{}

		if err := c.Bind().Body(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
