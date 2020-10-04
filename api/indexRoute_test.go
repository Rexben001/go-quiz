package index

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
)

func Test_index(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexRoute).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	apitest.New().
		Handler(r).
		Get("/").
		Expect(t).
		Status(http.StatusOK).
		End()
	apitest.New().
		Handler(r).
		Get("/").
		Expect(t).
		Body(`Welcome to Quiz's API`).
		End()
}

func Test_GetQuizzes(t *testing.T) {
	Index()
	r := mux.NewRouter()
	r.HandleFunc("/quizzes", GetALlQuizzes).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	apitest.New().
		Handler(r).
		Get("/quizzes").
		Expect(t).
		Status(http.StatusOK).
		End()
	// apitest.New().
	// 	Handler(r).
	// 	Get("/quizzes").
	// 	Expect(t).
	// 	Body(`{"message": "quiz fetched successfully"}`).
	// 	End()
}
