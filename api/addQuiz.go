package index

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddQuiz(response http.ResponseWriter, request *http.Request) {

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
	//json.NewDecoder() removes all but the Name field from each object
	json.NewDecoder(request.Body).Decode(&quiz)

	validateErr := validateAddQuiz(quiz)

	if validateErr != nil {
		response.WriteHeader(400)
		json.NewEncoder(response).Encode(validateErr)
		return
	}
	collection := getDB("quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, quiz)
	if err != nil {
		responseError(err, response)
		return
	}

	var id string

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		id = oid.Hex()
	}

	finalResult := createResult("New question added successfully", id)

	// writes the objects to standard output
	json.NewEncoder(response).Encode(finalResult)
}
