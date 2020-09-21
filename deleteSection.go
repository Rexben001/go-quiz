package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteSection(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	// get the params from the requst
	params := mux.Vars(request)

	database, _ := os.LookupEnv("DATABASE_NAME")
	secret, _ := os.LookupEnv("ACCESS_SECRET")

	response.Header().Add("content-type", "application/json")
	tokenString := request.Header.Get("Authorization")

	if string(tokenString) == "" {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Pls, provide a valid token"}`))
		return
	}
	updatedToken := strings.Split(tokenString, " ")[1]
	token, err := jwt.Parse(updatedToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Pls, provide a valid token"}`))
		return
	}

	// convert params id (string) to MongoDB ID
	id, _ := primitive.ObjectIDFromHex(params["id"])
	paramsID := string(params["id"])
	quizCollection := client.Database(database).Collection("quizzes")
	sectionCollection := client.Database(database).Collection("sections")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get item by id
	result, err := sectionCollection.DeleteOne(ctx, bson.M{"_id": id})
	resultQuiz, errQuiz := quizCollection.DeleteMany(ctx, bson.M{"owner": paramsID})

	if err != nil || errQuiz != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	if result.DeletedCount == 0 || resultQuiz.DeletedCount == 0 {
		// log.Fatal("Error on deleting one Hero", err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Unable to delete item"}`))
		return
	}
	finalResult := make(map[string]interface{})

	finalResult["message"] = "Section deleted successfully"
	finalResult["status"] = 200
	finalResult["success"] = true
	json.NewEncoder(response).Encode(finalResult)
}