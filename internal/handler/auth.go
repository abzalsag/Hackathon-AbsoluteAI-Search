package handler

import (
	"context"
	"encoding/json"
	"hackathon-absolutecinema/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthServiceInterface // Интерфейс
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON қатесі"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, `{"error":"Барлық өрістерді толтырыңыз"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	userID, err := h.authService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		http.Error(w, `{"error":"Тіркелу қатесі: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Тіркелу сәтті өтті",
		"user_id": userID,
	})
}

// Login POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON қатесі"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"Email және парольді енгізіңіз"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Кіру сәтті",
		"token":   token,
	})
}
