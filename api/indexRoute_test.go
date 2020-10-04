package index

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
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
	apitest.New().
		Handler(r).
		Get("/quizzes").
		Expect(t).
		Assert(jsonpath.Equal(`$.message`, "quiz fetched successfully")).
		End()
}
func Test_GetQuiz(t *testing.T) {
	Index()
	r := mux.NewRouter()
	r.HandleFunc("/quizzes/5f63d7ea28368128ce867c52", GetQuiz).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	apitest.New().
		Handler(r).
		Get("/quizzes/5f63d7ea28368128ce867c52").
		Expect(t).
		Status(http.StatusOK).
		End()
	apitest.New().
		Handler(r).
		Get("/quizzes/5f63d7ea28368128ce867c52").
		Expect(t).
		Assert(jsonpath.Equal(`$.message`, "Quiz fetched successfully")).
		End()
}
