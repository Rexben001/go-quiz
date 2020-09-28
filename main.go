package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	index "goQuiz/api"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	fmt.Println("App has started!!!!")

	port, _ := os.LookupEnv("PORT")

	index.Index()

	// define router
	router := mux.NewRouter()
	router.HandleFunc("/", index.IndexRoute)
	router.HandleFunc("/quizzes", index.AddQuiz).Methods("POST")
	router.HandleFunc("/quizzes", index.GetALlQuizzes).Methods("GET")

	router.HandleFunc("/quizzes/sections", index.AddSection).Methods("POST")
	router.HandleFunc("/quizzes/sections", index.GetAllSections).Methods("GET")

	router.HandleFunc("/quizzes/{id}", index.GetQuiz).Methods("GET")
	router.HandleFunc("/quizzes/{id}", index.UpdateQuiz).Methods("PUT")
	router.HandleFunc("/quizzes/{id}", index.DeleteQuiz).Methods("DELETE")

	router.HandleFunc("/quizzes/sections/{id}", index.GetQuizByOwner).Methods("GET")
	router.HandleFunc("/quizzes/sections/{id}", index.UpdateSection).Methods("PUT")
	router.HandleFunc("/quizzes/sections/{id}", index.DeleteSection).Methods("DELETE")

	router.HandleFunc("/quizzes/highscores", index.AddHighscore).Methods("POST")

	router.HandleFunc("/signup", index.CreateUser).Methods("POST")
	router.HandleFunc("/login", index.LoginUser).Methods("POST")

	http.ListenAndServe(":"+port, router)
}
