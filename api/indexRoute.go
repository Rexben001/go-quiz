package index

import (
	"net/http"
)

func IndexRoute(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte(`Welcome to Quiz's API`))

}
