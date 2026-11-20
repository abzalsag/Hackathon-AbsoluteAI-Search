package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"hackathon-absolutecinema/internal/handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, email, password, name string) (int64, error) {
	args := m.Called(ctx, email, password, name)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	mockService := new(MockAuthService)
	h := handler.NewAuthHandler(mockService)

	// RegisterRequest орнына тікелей JSON қолданамыз
	reqBody := `{"email":"test@test.com", "password":"123456", "name":"Test User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService.On("Register", mock.Anything, "test@test.com", "123456", "Test User").Return(int64(1), nil)

	h.Register(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login(t *testing.T) {
	mockService := new(MockAuthService)
	h := handler.NewAuthHandler(mockService)

	// LoginRequest орнына тікелей JSON қолданамыз
	reqBody := `{"email":"test@test.com", "password":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService.On("Login", mock.Anything, "test@test.com", "123456").Return("fake-jwt-token", nil)

	h.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}
