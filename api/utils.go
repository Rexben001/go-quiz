package index

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func getDB(name string) *mongo.Collection {

	database, _ := os.LookupEnv("DATABASE_NAME")
	collection := client.Database(database).Collection(name)

	return collection
}

func responseError(err error, response http.ResponseWriter) {
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte(`{"message": "` + err.Error() + `"}`))
}

func getResults(status int, message string, quizzes []Quizzes) map[string]interface{} {
	finalResult := make(map[string]interface{})
	finalResult["message"] = message
	finalResult["status"] = status
	finalResult["success"] = true
	finalResult["data"] = quizzes
	finalResult["totalQuizzes"] = len(quizzes)
	return finalResult
}

func getResultsSections(status int, message string, sections []Sections) map[string]interface{} {
	finalResult := make(map[string]interface{})
	finalResult["message"] = message
	finalResult["status"] = status
	finalResult["success"] = true
	finalResult["data"] = sections
	finalResult["totalSections"] = len(sections)
	return finalResult
}

func getResult(status int, message string, quiz Quizzes) map[string]interface{} {
	finalResult := make(map[string]interface{})
	finalResult["message"] = message
	finalResult["status"] = status
	finalResult["success"] = true
	finalResult["data"] = quiz
	return finalResult
}

func createResult(message string, id string) map[string]interface{} {

	finalResult := make(map[string]interface{})
	finalResult["message"] = message
	if id != "" {
		finalResult["InsertedId"] = id
	}
	finalResult["status"] = 201
	finalResult["success"] = true
	return finalResult

}

func validateToken(request *http.Request) (string, error) {

	secret, _ := os.LookupEnv("ACCESS_SECRET")

	tokenString := request.Header.Get("Authorization")

	if string(tokenString) == "" {
		return "", errors.New("Pls, provide a valid token")
	}

	updatedToken := strings.Split(tokenString, " ")[1]

	token, _ := jwt.Parse(updatedToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("Pls, provide a valid token")
		}
		return []byte(secret), nil
	})

	var errEmpty error
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["id"].(string), errEmpty
	}
	return "", errors.New("Pls, provide a valid token")

}
