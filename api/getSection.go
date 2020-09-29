package index

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllSections(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var sections []Sections
	collection := getDB("sections")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get all the items from the collection
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		responseError(err, response)
		return
	}

	defer cursor.Close(ctx)

	// iterate over the cursor and save the results as array
	for cursor.Next(ctx) {
		var section Sections
		cursor.Decode(&section)
		sections = append(sections, section)
	}
	// handle error
	if err := cursor.Err(); err != nil {
		responseError(err, response)
		return
	}

	finalResult := getResultsSections(200, "sections fetched successfully", sections)
	json.NewEncoder(response).Encode(finalResult)
}
