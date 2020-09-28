package index

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllSections(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var sections []Sections
	database, _ := os.LookupEnv("DATABASE_NAME")

	collection := client.Database(database).Collection("sections")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get all the items from the collection
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
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
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	finalResult := getResultsSections(200, "sections fetched successfully", sections)
	json.NewEncoder(response).Encode(finalResult)
}
