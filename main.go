package main

import (
	"net/http"
	"os"
	"task-manager-backend/handlers"
	"task-manager-backend/middleware"
	"task-manager-backend/repository"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
)


func main() {
	db := initDB()
	defer db.Close()
	
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	taskRepo := repository.NewTaskRepo(db) //создаёшь репозиторий задач, который работает с базой
	userRepo := repository.NewUserRepo(db) //создаёшь репозиторий пользователей

	taskHandler := &handlers.TaskHandler{
		Repo:    taskRepo,
		Log: log,
	}

	authHandler := &handlers.AuthHandler{
		Repo: userRepo,
		Log: log,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(log))
	// r.Use(middleware.MyMiddleware)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Get("/", taskHandler.GetTasks)
		r.Put("/{id}", taskHandler.UpdateTask)
	})


	log.Info().Msg("Server running on :8083")
	if err := http.ListenAndServe(":8083", r); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
