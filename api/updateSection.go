package index

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateSection(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var section Sections
	tokenID, err := validateToken(request)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Pls, provide a valid token"}`))
		return
	}

	section.UserID = tokenID

	// get the body request and decode it
	json.NewDecoder(request.Body).Decode(&section)
	// get the params from the requst
	params := mux.Vars(request)
	// convert params id (string) to MongoDB ID
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := getDB("sections")
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

	finalResult := createResult("Section updated successfully", "")

	// writes the objects to standard output
	json.NewEncoder(response).Encode(finalResult)
}
