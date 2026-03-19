package handlers

import (
	"encoding/json"
	"net/http"
	"task-manager-backend/repository"

	"github.com/rs/zerolog"
)

// это “объект, который умеет обрабатывать запросы по задачам”
type TaskHandler struct {
	Repo repository.TaskRepository
	Log  zerolog.Logger //Если в handler случится ошибка, он сможет записать её в лог.
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int) // раньше middleware/auth положил user_id в context
	tasks, err := h.Repo.GetAll(userID)
	if err != nil {
		h.Log.Error().Err(err).Msg("GetAll failed")                                 // Пишем ошибку в лог через zerolog
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError) //500
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
