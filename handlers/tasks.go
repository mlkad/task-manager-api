package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"task-manager-backend/models"
	"task-manager-backend/repository"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// это “объект, который умеет обрабатывать запросы по задачам”
type TaskHandler struct {
	Repo repository.TaskRepository
	Log  zerolog.Logger //Если в handler случится ошибка, он сможет записать её в лог.
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {     
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized) 
		return
	}

	doneStr := r.URL.Query().Get("done")
	priority := r.URL.Query().Get("priority")

	var done *bool
	if doneStr != "" {
		doneValue, err := strconv.ParseBool(doneStr)
		if err != nil {
			http.Error(w, `{"error":"invalid done value"}`, http.StatusBadRequest)
			return
		}
		done = &doneValue
	}

	if priority != "" && priority != "low" && priority != "medium" && priority != "high" {
		http.Error(w, `{"error":"invalid priority value"}`, http.StatusBadRequest)
		return
	} 
	
	var priorityFilter *string
	if priority != "" {
		priorityFilter = &priority
	}

	tasks, err := h.Repo.GetAll(userID, done, priorityFilter)
	if err != nil {
		h.Log.Error().Err(err).Msg("GetAll failed")
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	taskIDstr := chi.URLParam(r, "id") 
	taskID, err := strconv.Atoi(taskIDstr)
	if err != nil {
		http.Error(w, `{"error":"invalid task id"}`, http.StatusBadRequest)
		return
	}

	var input struct {
		Title string `json:"title"`
		Done bool `json:"done"`
		Priority string `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	task := models.Task{
		ID: taskID,
		Title: input.Title,
		Done: input.Done,
		Priority: input.Priority,
		UserID: userID,
	}

	err = h.Repo.Update(task)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound) //404
			return
		}
		
		h.Log.Error().Err(err).Msg("Update failed")
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

/*
func (h *AuthHandler) CreateUserWithTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Title    string `json:"title" validate:"required"`
		Done     bool   `json:"done"`
		Priority string `json:"priority" validate:"required,oneof=low medium high"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		http.Error(w, `{"error":"validation failed"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	task := models.Task{
		Title:    input.Title,
		Done:     input.Done,
		Priority: input.Priority,
	}

	err = h.Repo.CreateWithTask(user, task)
	if err != nil {
		h.Log.Error().Err(err).Msg("CreateWithTask failed")
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user and task created",
	})
}

*/