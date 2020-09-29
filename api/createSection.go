package index

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddSection(response http.ResponseWriter, request *http.Request) {

	var section Sections

	tokenID, err := validateToken(request)

	if err != nil {
		response.WriteHeader(400)
		response.Write([]byte(`{"message": "Pls, provide a valid token"}`))
		return
	}

	section.UserID = tokenID

	// get the body request and decode it
	//json.NewDecoder() removes all but the Name field from each object
	json.NewDecoder(request.Body).Decode(&section)

	collection := getDB("sections")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errTitle := collection.FindOne(ctx, bson.D{{"title", section.Title}}).Decode(&section)

	if errTitle == nil {
		response.WriteHeader(400)
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

	finalResult := createResult("New section added successfully", id)

	// writes the objects to standard output
	json.NewEncoder(response).Encode(finalResult)
}
