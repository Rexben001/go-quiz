package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var client *mongo.Client

func main() {
	fmt.Println("App has started!!!!")

	port, _ := os.LookupEnv("PORT")

	mongoURI, exists := os.LookupEnv("MONGO_URI")

	if exists {
		fmt.Println("ENV files loaded ")
	}

	// define timeout for Mongo and Go
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// mongodb connection
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	if client != nil {
		fmt.Println("Connected successfully")
	}

	// define router
	router := mux.NewRouter()
	router.HandleFunc("/", Index)

	http.ListenAndServe(":"+port, router)
}
