package index

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetQuiz(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var quiz Quizzes
	// get the params from the requst
	params := mux.Vars(request)
	// convert params id (string) to MongoDB ID
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := getDB("quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// get item by id
	err := collection.FindOne(ctx, Quizzes{ID: id}).Decode(&quiz)

	if err != nil {
		responseError(err, response)
		return
	}

	finalResult := getResult(200, "Quiz fetched successfully", quiz)
	json.NewEncoder(response).Encode(finalResult)
}
