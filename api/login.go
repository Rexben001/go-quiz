package index

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
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

	// if len(user.Email) < 4 {
	// 	newErr := errors.New("Invalid email or password >>input")
	// 	responseError(newErr, response)
	// 	return
	// }

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !re.MatchString(user.Email) {
		newErr := errors.New("Pls, enter a valid email address")
		responseError(newErr, response)
		return
	}

	collection := getDB("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// filter := bson.D{"email": user.Email}

	err := collection.FindOne(ctx, bson.D{{"email", user.Email}}).Decode(&result)
	// emailFound := collection.FindOne(ctx, Users{Email: user.Email})

	fmt.Println(err)

	if err != nil {
		newErr := errors.New("Invalid email or password")
		responseError(newErr, response)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		newErr := errors.New("Invalid email or password")
		responseError(newErr, response)
		return
	}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = result.ID
	atClaims["email"] = result.Email
	atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		newErr := errors.New("Unable to generate token")
		responseError(newErr, response)
		return
	}

	finalResult := make(map[string]interface{})

	finalResult["message"] = "User logged in successfully"
	finalResult["status"] = 200
	finalResult["success"] = true
	finalResult["token"] = token

	json.NewEncoder(response).Encode(finalResult)
}
