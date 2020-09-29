package index

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetQuizByOwner(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var quizzes []Quizzes

	json.NewDecoder(request.Body).Decode(&quizzes)

	collection := getDB("quizzes")
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

	finalResult := getResults(200, "quiz fetched by section successfully", quizzes)
	json.NewEncoder(response).Encode(finalResult)
}