package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID  string             `json:"userId,omitempty" bson:"userId,omitempty"`
	Caption string             `json:"caption,omitempty" bson:"caption,omitempty"`
	URL     string             `json:"url,omitempty" bson:"url,omitempty"`
	Time    string             `json:"time,omitempty" bson:"time,omitempty"`
}

var client *mongo.Client
var lock sync.Mutex

//endPoint - Create Post
func CreatePostEndpoint(response http.ResponseWriter, request *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	response.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)
	post.Time = time.Now().Format("2006-01-02 15:04:05")
	collection := client.Database("Instagram").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)

	time.Sleep(1 * time.Second)
}

//endPoint - Get Post
func GetPostEndpoint(response http.ResponseWriter, request *http.Request) {

	lock.Lock()
	defer lock.Unlock()

	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Post
	collection := client.Database("Instagram").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Post{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(post)

	time.Sleep(1 * time.Second)
}

//endPoint - Get User Post
func GetAllPostsofaUserEndpoint(response http.ResponseWriter, request *http.Request) {

	lock.Lock()
	defer lock.Unlock()

	var slicePost []Post
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := (params["id"])
	var post Post
	collection := client.Database("Instagram").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.M{"userId": id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	for cur.Next(ctx) {
		err := cur.Decode(&post)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}
		slicePost = append(slicePost, post)
	}
	for _, item := range slicePost {
		if item.UserID == id {
			fmt.Println(item.ID)
			fmt.Println(item.Caption)
			json.NewEncoder(response).Encode(item)
		}
	}

	time.Sleep(1 * time.Second)
}

func main() {
	fmt.Println("Starting the application")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()

	router.HandleFunc("/posts", CreatePostEndpoint).Methods("POST")
	router.HandleFunc("/post/{id}", GetPostEndpoint).Methods("GET")
	router.HandleFunc("/posts/users/{id}", GetAllPostsofaUserEndpoint).Methods("GET")

	http.ListenAndServe(":5000", router)
}
