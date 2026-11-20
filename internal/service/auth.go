package service

import (
	"context"
	"errors"
	"time"

	"hackathon-absolutecinema/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.UserRepository // Интерфейс
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (int64, error) {
	// Парольді хэштеу
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Дерекқорға жазу
	return s.userRepo.Create(ctx, email, string(hashedPassword), name)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Қолданушыны іздеу
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("қолданушы табылмады немесе пароль қате")
	}

	// Парольді тексеру
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("қолданушы табылмады немесе пароль қате")
	}

	// JWT токен жасау
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
