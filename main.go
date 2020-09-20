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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Quizzes struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Question string             `json:"question,omitempty" bson:"question,omitempty"`
	Options  []string           `json:"options,omitempty" bson:"options,omitempty"`
	Answer   string             `json:"answer,omitempty" bson:"answer,omitempty"`
	Owner    string             `json:"owner,omitempty" bson:"owner,omitempty"`
	UserID   string             `json:"userid,omitempty" bson:"userid,omitempty"`
}

type Users struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}
type Sections struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID string             `json:"userid,omitempty" bson:"userid,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
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
	router.HandleFunc("/quizzes", AddQuiz).Methods("POST")
	router.HandleFunc("/quizzes", GetALlQuizzes).Methods("GET")
	router.HandleFunc("/quizzes/{id}", GetQuiz).Methods("GET")
	router.HandleFunc("/quizzes/{id}", UpdateQuiz).Methods("PUT")
	router.HandleFunc("/quizzes/{id}", DeleteQuiz).Methods("DELETE")
	router.HandleFunc("/quizzes/section/{owner}", GetQuizByOwner).Methods("GET")
	router.HandleFunc("/quizzes/section", AddSection).Methods("POST")

	router.HandleFunc("/signup", CreateUser).Methods("POST")
	router.HandleFunc("/login", LoginUser).Methods("POST")

	http.ListenAndServe(":"+port, router)
}
