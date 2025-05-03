package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/hello", func(c fiber.Ctx) error {
		//Send a HTTP Status and a JSON File response to the client
		return c.Status(200).JSON(fiber.Map{"msg": "hello postman"})
	})

	app.Get("/teste", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Teste")
	})
	app.Get("/api/todos/", func(c fiber.Ctx) error {

		/*
			for i := range todos {
				fmt.Println(i)
			}
			fmt.Println("******************")
			for i, todo := range todos {
				fmt.Print(i, todo)
			}
		*/
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

		todo.ID = len(todos) + 1     // Start on 1
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

	// Delete a To-Do
	app.Delete("api/todos/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		// i is the position of array
		// todo is the object
		// todo.ID is the id of object
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				// Update todos list with all data before and after todos[i]
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"sucess": "true"})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Start the server on port...
	log.Fatal(app.Listen(":" + PORT))
}
