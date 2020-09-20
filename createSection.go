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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddSection(response http.ResponseWriter, request *http.Request) {

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
	//json.NewDecoder() removes all but the Name field from each object
	json.NewDecoder(request.Body).Decode(&section)

	collection := client.Database(database).Collection("sections")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errTitle := collection.FindOne(ctx, bson.D{{"title", section.Title}}).Decode(&section)

	if errTitle == nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Email already exists"}`))
		return
	}

	result, err := collection.InsertOne(ctx, section)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	var id string

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		id = oid.Hex()
	}

	finalResult := make(map[string]interface{})

	finalResult["message"] = "New question added successfully"
	finalResult["InsertedId"] = id
	finalResult["status"] = 201
	finalResult["success"] = true

	// writes the objects to standard output
	json.NewEncoder(response).Encode(finalResult)
}
