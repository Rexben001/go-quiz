package index

func validateAddQuiz(quiz Quizzes) map[string]interface{} {
	finalError := make(map[string]interface{})
	if len(quiz.Question) < 4 {
		finalError["Question"] = "Length of question must be greater than 4"
	}
	if len(quiz.Answer) == 0 {
		finalError["Answer"] = "Can't submit an empty answer"
	}

	if len(quiz.Owner) < 24 {
		finalError["Owner"] = "Enter a valid owner id"

	}
	return finalError
}
