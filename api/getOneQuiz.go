package index

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetQuiz(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	database, _ := os.LookupEnv("DATABASE_NAME")

	var quiz Quizzes
	// get the params from the requst
	params := mux.Vars(request)
	// convert params id (string) to MongoDB ID
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database(database).Collection("quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get item by id
	err := collection.FindOne(ctx, Quizzes{ID: id}).Decode(&quiz)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	finalResult := make(map[string]interface{})

	finalResult["message"] = "Quiz fetched successfully"
	finalResult["status"] = 200
	finalResult["success"] = true
	finalResult["data"] = quiz

	json.NewEncoder(response).Encode(finalResult)
}
