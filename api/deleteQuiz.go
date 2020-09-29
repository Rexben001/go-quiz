package index

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteQuiz(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	// get the params from the requst
	params := mux.Vars(request)

	response.Header().Add("content-type", "application/json")

	collection := getDB("quizzes")

	_, err := validateToken(request)

	if err != nil {
		response.WriteHeader(400)
		response.Write([]byte(`{"message": "Pls, provide a valid token"}`))
		return
	}

	// convert params id (string) to MongoDB ID
	id, _ := primitive.ObjectIDFromHex(params["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get item by id
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	if result.DeletedCount == 0 {
		// log.Fatal("Error on deleting one Hero", err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Unable to delete item"}`))
		return
	}
	finalResult := createResult("Quiz deleted successfully", "")

	// writes the objects to standard output
	json.NewEncoder(response).Encode(finalResult)
}
