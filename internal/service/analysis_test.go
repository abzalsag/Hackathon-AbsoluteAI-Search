package service_test

import (
	"context"
	"testing"
	"time"

	"hackathon-absolutecinema/internal/ai_client"
	"hackathon-absolutecinema/internal/models"
	"hackathon-absolutecinema/internal/repository/mocks"
	"hackathon-absolutecinema/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAIClient struct {
	mock.Mock
}

func (m *MockAIClient) Analyze(ctx context.Context, req ai_client.AnalysisRequest) (*ai_client.AnalysisResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ai_client.AnalysisResponse), args.Error(1)
}

func TestAnalysisService_CreateAnalysis(t *testing.T) {
	mockAnalysisRepo := new(mocks.MockAnalysisRepository)
	mockCompanyRepo := new(mocks.MockCompanyRepository)
	mockAIClient := new(MockAIClient)

	svc := service.NewAnalysisService(mockAnalysisRepo, mockCompanyRepo, mockAIClient)

	company := &models.Company{
		ID: 1, Name: "Test Co", Website: "test.com", Description: "Test",
	}

	mockCompanyRepo.On("GetByID", mock.Anything, int64(1)).Return(company, nil)
	mockAnalysisRepo.On("Create", mock.Anything, int64(1), int64(1), "pending").Return(int64(100), nil)
	mockAIClient.On("Analyze", mock.Anything, mock.AnythingOfType("ai_client.AnalysisRequest")).
		Return(&ai_client.AnalysisResponse{VisibilityScore: 85.0}, nil)
	mockAnalysisRepo.On("UpdateResult", mock.Anything, int64(100), mock.Anything, mock.Anything).Return(nil)
	mockAnalysisRepo.On("UpdateStatus", mock.Anything, int64(100), mock.Anything).Return(nil)

	analysisID, err := svc.CreateAnalysis(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(100), analysisID)

	time.Sleep(300 * time.Millisecond)
	mockCompanyRepo.AssertExpectations(t)
	mockAnalysisRepo.AssertExpectations(t)
	mockAIClient.AssertExpectations(t)
}
