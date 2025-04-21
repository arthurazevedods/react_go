package main

import (
	"fmt"
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
	app.Get("/todos", func(c fiber.Ctx) error {

		return c.Status(200).JSON(todos)
	})
	//Create a Todo
	app.Post("/api/todos/", func(c fiber.Ctx) error {
		todo := &Todo{} //Memory address of Todo struct

		if err := c.Bind().Body(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.ID = len(todos)
		todos = append(todos, *todo) // Pointer to the value stored at the memory address of todo

		return c.Status(201).JSON(todo)
	})

	// Update a Todo
	app.Patch("api/todos/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
