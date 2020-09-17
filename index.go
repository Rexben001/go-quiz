package main

import (
	"net/http"
)

func Index(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte(`Welcome to Quiz's API`))

}
