package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/arthurazevedods/react_go/config"
	"github.com/arthurazevedods/react_go/handlers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client // Declare mongoClient globally

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Inicializar cliente MongoDB usando config.ConnectDB()
	var err error
	mongoClient, err = config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Configurar rotas
	http.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTodos(mongoClient, w, r)
	})
	http.HandleFunc("/api/insertTodo", func(w http.ResponseWriter, r *http.Request) {
		handlers.InsertTodo(mongoClient, w, r)
	})
	http.HandleFunc("/api/insertMany", func(w http.ResponseWriter, r *http.Request) {
		handlers.InsertMany(mongoClient, w, r)
	})
	http.HandleFunc("/api/check", func(w http.ResponseWriter, r *http.Request) {
		handlers.CheckCollection(mongoClient, w, r)
	})
	http.HandleFunc("/api/delete", func(w http.ResponseWriter, r *http.Request) {
		handlers.Delete(mongoClient, w, r)
	})
	http.HandleFunc("/api/update", func(w http.ResponseWriter, r *http.Request) {
		handlers.Update(mongoClient, w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
