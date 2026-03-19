package handlers

import (
	"task-manager-backend/repository"

	"github.com/rs/zerolog"
)

type AuthHandler struct {
	Repo repository.UserRepository
	Log  zerolog.Logger
}
