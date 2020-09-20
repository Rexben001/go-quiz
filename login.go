package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	database, _ := os.LookupEnv("DATABASE_NAME")

	var user Users
	var result Users

	json.NewDecoder(request.Body).Decode(&user)

	collection := client.Database(database).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// filter := bson.D{"email": user.Email}

	err := collection.FindOne(ctx, bson.D{{"email", user.Email}}).Decode(&result)
	// emailFound := collection.FindOne(ctx, Users{Email: user.Email})

	fmt.Println(err)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Email or password is incorrect"}`))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		fmt.Println("Wrong password")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Email or password is incorrect"}`))
		return
	}

	finalResult := make(map[string]interface{})

	finalResult["message"] = "User logged in successfully"
	finalResult["status"] = 200
	finalResult["success"] = true

	json.NewEncoder(response).Encode(finalResult)
}
