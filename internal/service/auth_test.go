package service_test

import (
	"context"
	"testing"

	"hackathon-absolutecinema/internal/models"
	"hackathon-absolutecinema/internal/repository/mocks"
	"hackathon-absolutecinema/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	t.Run("Сәтті тіркелу", func(t *testing.T) {
		// ⚠ Әр subtest-те ЖАҢА mock
		mockUserRepo := new(mocks.MockUserRepository)
		svc := service.NewAuthService(mockUserRepo, "secret_key")

		mockUserRepo.On("Create", mock.Anything, "test@test.com", mock.AnythingOfType("string"), "USER").
			Return(int64(1), nil).
			Once()

		_, err := svc.Register(context.Background(), "test@test.com", "password123", "USER")

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Database қатесі", func(t *testing.T) {
		// ⚠ Әр subtest-те ЖАҢА mock — ескі nil expectation жоқ
		mockUserRepo := new(mocks.MockUserRepository)
		svc := service.NewAuthService(mockUserRepo, "secret_key")

		mockUserRepo.On("Create", mock.Anything, "test@test.com", mock.AnythingOfType("string"), "USER").
			Return(int64(0), assert.AnError).
			Once()

		_, err := svc.Register(context.Background(), "test@test.com", "password123", "USER")

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	plainPassword := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	assert.NoError(t, err)

	newUser := func() *models.User {
		return &models.User{
			ID:           1,
			Email:        "test@test.com",
			PasswordHash: string(hashedPassword),
		}
	}

	t.Run("Сәтті кіру", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		svc := service.NewAuthService(mockUserRepo, "secret_key")

		mockUserRepo.On("GetByEmail", mock.Anything, "test@test.com").
			Return(newUser(), nil).
			Once()

		token, err := svc.Login(context.Background(), "test@test.com", plainPassword)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Қолданушы табылмады", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		svc := service.NewAuthService(mockUserRepo, "secret_key")

		mockUserRepo.On("GetByEmail", mock.Anything, "notfound@test.com").
			Return(nil, assert.AnError).
			Once()

		_, err := svc.Login(context.Background(), "notfound@test.com", plainPassword)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Қате пароль", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		svc := service.NewAuthService(mockUserRepo, "secret_key")

		mockUserRepo.On("GetByEmail", mock.Anything, "test@test.com").
			Return(newUser(), nil).
			Once()

		_, err := svc.Login(context.Background(), "test@test.com", "wrong_password")

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}
