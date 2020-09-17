package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetALlQuizzes(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var quizzes []Quizzes
	database, _ := os.LookupEnv("DATABASE_NAME")

	collection := client.Database(database).Collection("quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get all the items from the collection
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
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
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	finalResult := make(map[string]interface{})

	finalResult["message"] = "quiz fetched successfully"
	finalResult["status"] = 200
	finalResult["success"] = true
	finalResult["data"] = quizzes
	finalResult["totalQuizzes"] = len(quizzes)
	json.NewEncoder(response).Encode(finalResult)
}
