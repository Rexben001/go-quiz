package index

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

func UpdateSection(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

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

	var section Sections

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		section.UserID = claims["id"].(string)
	} else {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Pls, provide a valid token"}`))
		return
	}

	// get the body request and decode it
	json.NewDecoder(request.Body).Decode(&section)
	// get the params from the requst
	params := mux.Vars(request)
	// convert params id (string) to MongoDB ID
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database(database).Collection("quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get item by id
	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": section})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	if result.ModifiedCount == 0 {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Unable to update item"}`))
		return
	}
	finalResult := make(map[string]interface{})
	finalResult["message"] = "Section updated successfully"
	finalResult["status"] = 201
	finalResult["success"] = true
	json.NewEncoder(response).Encode(finalResult)
}
