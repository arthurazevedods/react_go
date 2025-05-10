package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Todo struct {
	Completed bool   `json:"completed" bson:"completed"`
	Body      string `json:"body" bson:"body"`
}

func main() {
	fmt.Println("Hello, World!")
	var uri string
	// Read .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.")
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("react_go").Collection("Todos")
	Todos := []Todo{
		{Completed: false, Body: "Buy groceries"},
		{Completed: false, Body: "Go to the gym"},
		{Completed: false, Body: "Read a book"},
		{Completed: false, Body: "Learn Go"},
		{Completed: false, Body: "Learn React"},
		{Completed: false, Body: "Learn MongoDB"},
		{Completed: false, Body: "Learn Docker"},
		{Completed: false, Body: "Learn Kubernetes"},
		{Completed: false, Body: "Learn GraphQL"},
		{Completed: false, Body: "Learn TypeScript"},
	}
	result, err := collection.InsertMany(context.TODO(), Todos)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted many documents: ", result.InsertedIDs)

}
