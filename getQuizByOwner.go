package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetQuizByOwner(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	database, _ := os.LookupEnv("DATABASE_NAME")

	var quizzes []Quizzes

	json.NewDecoder(request.Body).Decode(&quizzes)

	collection := client.Database(database).Collection("quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := mux.Vars(request)
	owner := string(params["id"])

	// fmt.Println("params", params)
	fmt.Println("owner", owner)

	cursor, err := collection.Find(ctx, bson.D{{"owner", owner}})
	// emailFound := collection.FindOne(ctx, Users{Email: user.Email})

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Email or password is incorrect"}`))
		return
	}

	defer cursor.Close(ctx)

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

	finalResult["message"] = "quizzes by owner fetched successfully"
	finalResult["status"] = 200
	finalResult["success"] = true
	finalResult["data"] = quizzes
	finalResult["totalQuizzes"] = len(quizzes)
	json.NewEncoder(response).Encode(finalResult)
}
