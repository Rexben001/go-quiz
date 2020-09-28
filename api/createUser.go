package index

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(response http.ResponseWriter, request *http.Request) {

	database, _ := os.LookupEnv("DATABASE_NAME")

	response.Header().Add("content-type", "application/json")
	var user Users

	// get the body request and decode it
	//json.NewDecoder() removes all but the Name field from each object
	json.NewDecoder(request.Body).Decode(&user)

	collection := client.Database(database).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errEmail := collection.FindOne(ctx, bson.D{{"email", user.Email}}).Decode(&user)

	if errEmail == nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Email already exists"}`))
		return
	}

	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.Password = string(hashedPassword)
	if errPassword != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Unable to create an account. Try again later"}`))
		return
	}

	result, err := collection.InsertOne(ctx, user)
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
