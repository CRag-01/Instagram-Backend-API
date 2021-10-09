/* So basically what i understood by pagination is, limiting the fetch(GET requests) in case of any networks issues or such
   How I went about pagination here is, limiting the load of posts by limiting the the no of posts of a single use,
   NOTE: LIMIT PARAMETER IS USED IN THE URL WHILE ROUTING
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

//PAGINATION - END POINT
func PaginationEndPoint(response http.ResponseWriter, request *http.Request) {

	lock.Lock()
	defer lock.Unlock()

	var slicePost []Post
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := (params["id"])
	limit, _ := strconv.Atoi(params["limit"])

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
			if limit > 0 {
				limit--
				json.NewEncoder(response).Encode(item)
			}
		}
	}

	time.Sleep(1 * time.Second)
}

func main() {
	fmt.Println("Starting the application")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()

	router.HandleFunc("/posts/users/{id}&SetLimit={limit}", PaginationEndPoint).Methods("GET")

	http.ListenAndServe(":5000", router)
}
