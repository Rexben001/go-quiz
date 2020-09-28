package index

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
type Highscores struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User    string             `json:"user,omitempty" bson:"user,omitempty"`
	Section string             `json:"section,omitempty" bson:"section,omitempty"`
	Score   int                `json:"score,omitempty" bson:"score,omitempty"`
}

var client *mongo.Client

func Index() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
}
