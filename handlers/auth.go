package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"task-manager-backend/models"
	"task-manager-backend/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Repo repository.UserRepository
	Log  zerolog.Logger
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, `{"error":"validation failed"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	user, err := h.Repo.Create(models.User{Email: req.Email, PasswordHash: string(hash) })
	if err != nil {
		http.Error(w, `{"error":"email already exists"}`, http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": user.ID})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string`json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req) 

	user, err := h.Repo.FindByEmail(req.Email)

	//не говорим ЧТО именно не так — email или пароль
  if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
    http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
    return
  }

	//Access token — короткоживущий (15 минут)
	accessToken, _ := generateToken(user.ID, 15*time.Minute)
	
	//refresh - долгоживущий (неделя)
	refreshToken, _ := generateToken(user.ID, 7*24*time.Hour)

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func generateToken(userID int, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}


func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, `{"error":"invalid refresh token}`, http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))
	
	newAccess, _ := generateToken(userID, 15*time.Minute)
	json.NewEncoder(w).Encode(map[string]string{"access_token": newAccess})
}