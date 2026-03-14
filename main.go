package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Task struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Done bool `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	Password string `json:"-"`
	Priority string `json:"priority"`
}

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required,min=3"`
	Priority string `json:"priority"`
}

var tasks = []Task{
	{ID: 1, Title: "learn golang", Done: false, CreatedAt: time.Now(), Priority: "low"},
	{ID: 2, Title: "buy milk", Done: true, Priority: "low"},
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler)

	r.Post("/tasks", createTask)
	r.Get("/tasks", getTasks)

	r.Get("/tasks/{id}", getTask)
	r.Put("/tasks/{id}", updateTask)
	r.Delete("/tasks/{id}", deleteTask)
	
	log.Println("Server running on :8085")
	log.Fatal(http.ListenAndServe(":8085", r))
} 

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, I am your backend")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

var validate = validator.New()

func createTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return 
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest) //400
		return
	}

	priority := req.Priority

	if priority == "" {
		priority = "low"
	}

	task := Task{
		ID: len(tasks) + 1,
		Title: req.Title,
		Done: false,
		CreatedAt: time.Now(),
		Priority: priority,
	}

	tasks = append(tasks, task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
} 

func updateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	for i, t := range tasks {
		if t.ID == idInt {
			tasks[i].Title = req.Title
			if req.Priority != "" {
				tasks[i].Priority = req.Priority
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.Error(w, "task not found", http.StatusNotFound) //404
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for i, t := range tasks {
		if t.ID == idInt {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent) //delete 204
			return
		}
	}
	http.Error(w, "task not found", http.StatusNotFound)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	for _, t := range tasks {
		if idInt == t.ID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

/*
package main

import (
 "context"
 "encoding/json"
 "log"
 "net/http"

 "github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

type User struct {
 ID   int    json:"id"
 Name string json:"name"
}

func main() {
 var err error
 conn, err = pgx.Connect(context.Background(),
  "postgres://postgres:password@localhost:5432/mydb")
 if err != nil {
  log.Fatal(err)
 }

 http.HandleFunc("/users", usersHandler)

 log.Println("Server started on :8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
 switch r.Method {

 case http.MethodGet:
  getUsers(w)

 case http.MethodPost:
  createUser(w, r)

 default:
  http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
 }
}

func getUsers(w http.ResponseWriter) {
 rows, err := conn.Query(context.Background(),
  "SELECT id, name FROM users")
 if err != nil {
  http.Error(w, err.Error(), 500)
  return
 }
 defer rows.Close()

 var users []User

 for rows.Next() {
  var u User
  err := rows.Scan(&u.ID, &u.Name)
  if err != nil {
   http.Error(w, err.Error(), 500)
   return
  }
  users = append(users, u)
 }

 json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
 var u User

 err := json.NewDecoder(r.Body).Decode(&u)
 if err != nil {
  http.Error(w, err.Error(), 400)
  return
 }

 err = conn.QueryRow(context.Background(),
  "INSERT INTO users (name) VALUES ($1) RETURNING id",
  u.Name).Scan(&u.ID)
 if err != nil {
  http.Error(w, err.Error(), 500)
  return
 }

 w.WriteHeader(http.StatusCreated)
 json.NewEncoder(w).Encode(u)
}
*/
