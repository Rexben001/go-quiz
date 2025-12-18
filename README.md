# go-quiz

REST API for managing quiz sections, questions, users, and highscores. Built with Go, Gorilla Mux, JWT auth, and MongoDB.

## Requirements
- Go 1.15+
- MongoDB instance (local or hosted)

## Environment
Create a `.env` file in the project root:

```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
DATABASE_NAME=go_quiz
ACCESS_SECRET=replace-with-a-strong-secret
```

## Run locally
```bash
go run .
```

The API listens on `http://localhost:<PORT>`.

## Authentication
Login returns a JWT. For protected endpoints, include:

```
Authorization: Bearer <token>
```

Protected endpoints include creating/updating/deleting quizzes and sections, and deleting quizzes/sections.

## Data models (JSON)

Quiz:
```json
{
  "question": "What is Go?",
  "options": ["Language", "Animal", "Both"],
  "answer": "Both",
  "owner": "603d2f5b64e2f2b2a9d2e0f1"
}
```

Section:
```json
{
  "title": "Go Basics"
}
```

User:
```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

Highscore:
```json
{
  "user": "user@example.com",
  "section": "Go Basics",
  "score": 9
}
```

Notes:
- `owner` is the section ID (24-char hex string).
- `userid` is injected server-side from the JWT token.
- Quiz validation requires a question length >= 4, a non-empty answer, and a valid owner ID.

## Endpoints

### Health
- `GET /` -> `Welcome to Quiz's API`

### Users
- `POST /signup` -> create user
- `POST /login` -> login and return token

### Quizzes
- `GET /quizzes` -> list all quizzes
- `POST /quizzes` -> create quiz (auth required)
- `GET /quizzes/{id}` -> get quiz by ID
- `PUT /quizzes/{id}` -> update quiz (auth required)
- `DELETE /quizzes/{id}` -> delete quiz (auth required)

### Sections
- `GET /quizzes/sections` -> list all sections
- `POST /quizzes/sections` -> create section (auth required)
- `GET /quizzes/sections/{id}` -> list quizzes by section (owner id)
- `PUT /quizzes/sections/{id}` -> update section (auth required)
- `DELETE /quizzes/sections/{id}` -> delete section and its quizzes (auth required)

### Highscores
- `POST /quizzes/highscores` -> add a highscore

## Example requests

Create user:
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'
```

Login:
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'
```

Create section:
```bash
curl -X POST http://localhost:8080/quizzes/sections \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Go Basics"}'
```

Create quiz:
```bash
curl -X POST http://localhost:8080/quizzes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"question":"What is Go?","options":["Language","Animal","Both"],"answer":"Both","owner":"603d2f5b64e2f2b2a9d2e0f1"}'
```

List quizzes by section:
```bash
curl http://localhost:8080/quizzes/sections/603d2f5b64e2f2b2a9d2e0f1
```

Add highscore:
```bash
curl -X POST http://localhost:8080/quizzes/highscores \
  -H "Content-Type: application/json" \
  -d '{"user":"user@example.com","section":"Go Basics","score":9}'
```

## Responses
Most success responses include:
```json
{
  "message": "string",
  "status": 200,
  "success": true
}
```

List responses include `data` and totals:
```json
{
  "message": "quiz fetched successfully",
  "status": 200,
  "success": true,
  "data": [],
  "totalQuizzes": 0
}
```

Errors are returned with status 400/404 and:
```json
{"message": "error text"}
```

## Tests
```bash
go test ./...
```

Some tests expect a reachable MongoDB configured by `.env`.
