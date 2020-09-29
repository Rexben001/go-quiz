package index

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	secret, _ := os.LookupEnv("ACCESS_SECRET")

	var user Users
	var result Users

	json.NewDecoder(request.Body).Decode(&user)

	collection := getDB("users")
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

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = result.ID
	atClaims["email"] = result.Email
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Unable to create token"}`))
		return
	}

	finalResult := make(map[string]interface{})

	finalResult["message"] = "User logged in successfully"
	finalResult["status"] = 200
	finalResult["success"] = true
	finalResult["token"] = token

	json.NewEncoder(response).Encode(finalResult)
}
