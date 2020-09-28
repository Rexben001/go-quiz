package index

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetALlQuizzes(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var quizzes []Quizzes

	collection := getDB("quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// get all the items from the collection
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		// response.WriteHeader(http.StatusInternalServerError)
		// response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		// return
		responseError(err, response)
		return
	}

	defer cursor.Close(ctx)

	// iterate over the cursor and save the results as array
	for cursor.Next(ctx) {
		var quiz Quizzes
		cursor.Decode(&quiz)
		quizzes = append(quizzes, quiz)
	}
	// handle error
	if err := cursor.Err(); err != nil {
		// response.WriteHeader(http.StatusInternalServerError)
		// response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		responseError(err, response)
		return
	}

	finalResult := getResults(200, "quiz fetched successfully", quizzes)
	json.NewEncoder(response).Encode(finalResult)
}
