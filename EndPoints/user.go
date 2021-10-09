package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	// L "./lib"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

var client *mongo.Client
var lock sync.Mutex

//endpoint - CREATE USER
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {

	lock.Lock()
	defer lock.Unlock()

	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	// fmt.Println(user.Password)
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hai bro")
	}
	user.Password = string(hash)
	collection := client.Database("Instagram").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)

	time.Sleep(1 * time.Second)
}

//endpoint - GET USER
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {

	lock.Lock()
	defer lock.Unlock()

	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("Instagram").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)

	time.Sleep(1 * time.Second)
}

func main() {
	fmt.Println("Starting the application")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()

	router.HandleFunc("/users", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/user/{id}", GetUserEndpoint).Methods("GET")

	http.ListenAndServe(":5000", router)

	// L.CheckPass()
}
