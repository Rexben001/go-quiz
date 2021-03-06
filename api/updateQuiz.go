package index

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateQuiz(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var quiz Quizzes
	tokenID, err := validateToken(request)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Pls, provide a valid token"}`))
		return
	}
	quiz.UserID = tokenID
	// get the body request and decode it
	json.NewDecoder(request.Body).Decode(&quiz)
	// get the params from the requst
	params := mux.Vars(request)
	// convert params id (string) to MongoDB ID
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := getDB("quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// get item by id
	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": quiz})

	if err != nil {
		responseError(err, response)
		return
	}

	if result.ModifiedCount == 0 {
		newErr := errors.New("Unable to update items")
		responseError(newErr, response)
		return
	}

	finalResult := createResult("Quiz updated successfully", "")

	// writes the objects to standard output
	json.NewEncoder(response).Encode(finalResult)
}
