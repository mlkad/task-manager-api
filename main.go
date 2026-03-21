package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"task-manager-backend/middleware"
	"time"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"-"`
	Priority  string    `json:"priority"`
}

type TaskHandler struct {
	tasks    []Task
	validate *validator.Validate
}

type CreateTaskRequest struct {
	Title    string `json:"title" validate:"required,min=3"`
	Priority string `json:"priority"`
}

// var tasks = []Task{
// 	{ID: 1, Title: "learn golang", Done: false, CreatedAt: time.Now(), Priority: "low"},
// 	{ID: 2, Title: "buy milk", Done: true, Priority: "low"},
// }

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	handler := &TaskHandler{
		tasks:    []Task{},
		validate: validator.New(),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(log))
	// r.Use(middleware.MyMiddleware)

	r.Get("/", handler.homeHandler)
	r.Get("/health", handler.healthHandler)

	r.Route("/tasks", func(r chi.Router) {
		r.Use(middleware.APIKeyMiddleware)
		r.Post("/", handler.createTask)
		r.Get("/", handler.getTasks)

		r.Get("/{id}", handler.getTask)
		r.Put("/{id}", handler.updateTask)
		r.Delete("/{id}", handler.deleteTask)
	})

	log.Info().Msg("Server running on :8083")
	if err := http.ListenAndServe(":8083", r); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func (h *TaskHandler) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, I am your backend")
}

func (h *TaskHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	err = h.validate.Struct(req)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest) //400
		return
	}

	priority := req.Priority

	if priority == "" {
		priority = "low"
	}

	task := Task{
		ID:        len(h.tasks) + 1,
		Title:     req.Title,
		Done:      false,
		CreatedAt: time.Now(),
		Priority:  priority,
	}

	h.tasks = append(h.tasks, task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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

	err = h.validate.Struct(req)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	for i, t := range h.tasks {
		if t.ID == idInt {
			h.tasks[i].Title = req.Title
			if req.Priority != "" {
				h.tasks[i].Priority = req.Priority
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(h.tasks[i])
			return
		}
	}
	http.Error(w, "task not found", http.StatusNotFound) //404
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for i, t := range h.tasks {
		if t.ID == idInt {
			h.tasks = append(h.tasks[:i], h.tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent) //delete 204
			return
		}
	}
	http.Error(w, "task not found", http.StatusNotFound)
}

func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.tasks)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	for _, t := range h.tasks {
		if idInt == t.ID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}
