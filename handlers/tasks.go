package handlers

import (
	"encoding/json"
	"net/http"
	"task-manager-backend/repository"

	"github.com/rs/zerolog"
)

type TaskHandler struct {
	Repo repository.TaskRepository
	Log zerolog.Logger
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	tasks, err := h.Repo.GetAll(userID)
	if err != nil {
		h.Log.Error().Err(err).Msg("GetAll failed")
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError) //500
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}