package service_test

import (
	"context"
	"testing"

	"hackathon-absolutecinema/internal/models"
	"hackathon-absolutecinema/internal/repository/mocks"
	"hackathon-absolutecinema/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompanyService_CreateCompany(t *testing.T) {
	t.Run("Сәтті компания қосу", func(t *testing.T) {
		mockCompanyRepo := new(mocks.MockCompanyRepository)
		svc := service.NewCompanyService(mockCompanyRepo)

		mockCompanyRepo.On("Create", mock.Anything, int64(1), "Test Co", "test.com", "Test Description").
			Return(int64(1), nil).
			Once()

		id, err := svc.CreateCompany(context.Background(), 1, "Test Co", "test.com", "Test Description")

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		mockCompanyRepo.AssertExpectations(t)
	})

	t.Run("Database қатесі", func(t *testing.T) {
		mockCompanyRepo := new(mocks.MockCompanyRepository)
		svc := service.NewCompanyService(mockCompanyRepo)

		mockCompanyRepo.On("Create", mock.Anything, int64(1), "Test Co", "test.com", "Test Description").
			Return(int64(0), assert.AnError).
			Once()

		_, err := svc.CreateCompany(context.Background(), 1, "Test Co", "test.com", "Test Description")

		assert.Error(t, err)
		mockCompanyRepo.AssertExpectations(t)
	})
}

func TestCompanyService_GetCompanies(t *testing.T) {
	mockCompanyRepo := new(mocks.MockCompanyRepository)
	svc := service.NewCompanyService(mockCompanyRepo)

	companies := []models.Company{
		{ID: 1, Name: "Test Co"},
	}

	mockCompanyRepo.On("GetByUserID", mock.Anything, int64(1)).
		Return(companies, nil).
		Once()

	result, err := svc.GetCompanies(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockCompanyRepo.AssertExpectations(t)
}

func TestCompanyService_UpdateCompany(t *testing.T) {
	mockCompanyRepo := new(mocks.MockCompanyRepository)
	svc := service.NewCompanyService(mockCompanyRepo)

	mockCompanyRepo.On("Update", mock.Anything, int64(1), "Updated Co", "updated.com", "Updated").
		Return(nil).
		Once()

	err := svc.UpdateCompany(context.Background(), 1, "Updated Co", "updated.com", "Updated")

	assert.NoError(t, err)
	mockCompanyRepo.AssertExpectations(t)
}

func TestCompanyService_DeleteCompany(t *testing.T) {
	mockCompanyRepo := new(mocks.MockCompanyRepository)
	svc := service.NewCompanyService(mockCompanyRepo)

	mockCompanyRepo.On("Delete", mock.Anything, int64(1)).
		Return(nil).
		Once()

	err := svc.DeleteCompany(context.Background(), 1)

	assert.NoError(t, err)
	mockCompanyRepo.AssertExpectations(t)
}
