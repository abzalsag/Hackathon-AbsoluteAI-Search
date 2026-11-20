package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"hackathon-absolutecinema/internal/handler"
	"hackathon-absolutecinema/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Company Service
type MockCompanyService struct {
	mock.Mock
}

func (m *MockCompanyService) GetCompanies(ctx context.Context, userID int64) ([]models.Company, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Company), args.Error(1)
}

func (m *MockCompanyService) CreateCompany(ctx context.Context, userID int64, name, website, description string) (int64, error) {
	args := m.Called(ctx, userID, name, website, description)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCompanyService) GetCompany(ctx context.Context, id int64) (*models.Company, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyService) UpdateCompany(ctx context.Context, id int64, name, website, description string) error {
	args := m.Called(ctx, id, name, website, description)
	return args.Error(0)
}

// DeleteCompany методын қостық!
func (m *MockCompanyService) DeleteCompany(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCompanyHandler_CreateCompany(t *testing.T) {
	mockService := new(MockCompanyService)
	h := handler.NewCompanyHandler(mockService)

	// CreateCompanyRequest орнына тікелей JSON қолданамыз
	reqBody := `{"name":"My Startup", "website":"https://test.com", "description":"AI startup"}`
	req := httptest.NewRequest(http.MethodPost, "/api/companies", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Middleware-тен келетін user_id-ні контекстке қосамыз
	ctx := context.WithValue(req.Context(), "user_id", int64(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	mockService.On("CreateCompany", mock.Anything, int64(1), "My Startup", "https://test.com", "AI startup").Return(int64(1), nil)

	h.CreateCompany(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}
