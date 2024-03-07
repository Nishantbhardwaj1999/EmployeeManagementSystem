package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Employee struct
type Employee struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Position string             `json:"position,omitempty" bson:"position,omitempty"`
	Age      int                `json:"age,omitempty" bson:"age,omitempty"`
}

var client *mongo.Client
var collection *mongo.Collection

func main() {
	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(context.Background(), clientOptions)
	collection = client.Database("employee_management").Collection("employees")

	// Initialize router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/employees", GetEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", GetEmployee).Methods("GET")
	router.HandleFunc("/employees", CreateEmployee).Methods("POST")
	router.HandleFunc("/employees/{id}", UpdateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", DeleteEmployee).Methods("DELETE")

	// Enable CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	// Start server with CORS middleware
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

// GetEmployees returns all employees
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	var employees []Employee
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var emp Employee
		cur.Decode(&emp)
		employees = append(employees, emp)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(employees)
}

// GetEmployee returns a specific employee
func GetEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var emp Employee
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&emp)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(emp)
}

// CreateEmployee creates a new employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	_ = json.NewDecoder(r.Body).Decode(&emp)
	_, err := collection.InsertOne(context.Background(), emp)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(emp)
}

// UpdateEmployee updates an existing employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var emp Employee
	_ = json.NewDecoder(r.Body).Decode(&emp)
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": emp})
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(emp)
}

// DeleteEmployee deletes an employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)
}
