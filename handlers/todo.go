package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/arthurazevedods/react_go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	Completed bool   `json:"completed" bson:"completed"`
	Body      string `json:"body" bson:"body"`
}

func GetTodos(mongoClient *mongo.Client, w http.ResponseWriter, r *http.Request) {
	collection := mongoClient.Database("react_go").Collection("Todos")
	curr, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		log.Printf("Failed to find documents: %v", err)
		http.Error(w, "Failed to find documents", http.StatusInternalServerError)
		return
	}
	defer curr.Close(r.Context())

	var todos []models.Todo
	for curr.Next(r.Context()) {
		var todo models.Todo
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

func InsertTodo(mongoClient *mongo.Client, w http.ResponseWriter, r *http.Request) {
	collection := mongoClient.Database("react_go").Collection("Todos")
	var newTodo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	result, err := collection.InsertOne(context.TODO(), newTodo)
	if err != nil {
		log.Printf("Failed to insert document: %v", err)
		http.Error(w, "Failed to insert document", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result.InsertedID); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func InsertMany(mongoClient *mongo.Client, w http.ResponseWriter, r *http.Request) {
	collection := mongoClient.Database("react_go").Collection("Todos")
	var todos []models.Todo
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

func CheckCollection(mongoClient *mongo.Client, w http.ResponseWriter, r *http.Request) {

	count, err := mongoClient.Database("react_go").Collection("Todos").CountDocuments(r.Context(), bson.D{})
	if err != nil {
		log.Printf("Failed to count documents: %v", err)
		http.Error(w, "Failed to count documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Total documents in collection: %d", count)
}

func Delete(mongoClient *mongo.Client, w http.ResponseWriter, r *http.Request) {
	// Obter o valor do parâmetro "id" da URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	// Converter o ID para o tipo ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Criar o filtro para localizar o documento
	filter := bson.D{{Key: "_id", Value: objectID}}

	// Obter a coleção e deletar o documento
	collection := mongoClient.Database("react_go").Collection("Todos")
	log.Println("Deleting document with ID:", id)
	log.Println("Filter:", filter)
	log.Println("Context:", r.Context())
	result := collection.FindOneAndDelete(r.Context(), filter)

	// Verificar se houve erro ao deletar
	if err := result.Err(); err != nil {
		log.Printf("Failed to delete document: %v", err)
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	// Decodificar o documento deletado
	var deletedTodo models.Todo
	if err := result.Decode(&deletedTodo); err != nil {
		log.Printf("Failed to decode deleted document: %v", err)
		http.Error(w, "Failed to decode deleted document", http.StatusInternalServerError)
		return
	}

	// Retornar o documento deletado como resposta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(deletedTodo); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Logar e enviar resposta de sucesso
	log.Printf("Deleted document: %v", deletedTodo)
	w.WriteHeader(http.StatusOK)
}
