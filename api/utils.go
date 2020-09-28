package index

import (
	"net/http"
	"os"

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
