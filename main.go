package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arthurazevedods/react_go/config"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	Completed bool   `json:"completed" bson:"completed"`
	Body      string `json:"body" bson:"body"`
}

var mongoClient *mongo.Client

func GetTodos(w http.ResponseWriter, r *http.Request) {
	collection := mongoClient.Database("react_go").Collection("Todos")
	curr, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		log.Printf("Failed to find documents: %v", err)
		http.Error(w, "Failed to find documents", http.StatusInternalServerError)
		return
	}
	defer curr.Close(r.Context())

	var todos []Todo
	for curr.Next(r.Context()) {
		var todo Todo
		if err := curr.Decode(&todo); err != nil {
			log.Printf("Failed to decode document: %v", err)
			http.Error(w, "Failed to decode document", http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	if err := curr.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		http.Error(w, "Cursor error", http.StatusInternalServerError)
		return
	}

	if len(todos) == 0 {
		log.Println("No todos found in collection")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func InsertMany(w http.ResponseWriter, r *http.Request) {
	collection := mongoClient.Database("react_go").Collection("Todos")
	var todos []Todo
	if err := json.NewDecoder(r.Body).Decode(&todos); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	docs := make([]interface{}, len(todos))
	for i, todo := range todos {
		docs[i] = todo
	}

	result, err := collection.InsertMany(r.Context(), docs)
	if err != nil {
		log.Printf("Failed to insert documents: %v", err)
		http.Error(w, "Failed to insert documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result.InsertedIDs); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CheckCollection(w http.ResponseWriter, r *http.Request) {
	count, err := mongoClient.Database("react_go").Collection("Todos").CountDocuments(r.Context(), bson.D{})
	if err != nil {
		log.Printf("Failed to count documents: %v", err)
		http.Error(w, "Failed to count documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Total documents in collection: %d", count)
}

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
	http.HandleFunc("/api", GetTodos)
	http.HandleFunc("/api/insertMany", InsertMany)
	http.HandleFunc("/api/check", CheckCollection)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
