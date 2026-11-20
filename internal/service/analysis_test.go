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

// Mock AI Client
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

	// 1. Компанияны табуды mock қылу
	mockCompanyRepo.On("GetByID", mock.Anything, int64(1)).Return(company, nil)

	// 2. Анализ жазбасын жасауды mock қылу (ID = 100)
	mockAnalysisRepo.On("Create", mock.Anything, int64(1), int64(1), "pending").Return(int64(100), nil)

	// 3. AI шақырылымын mock қылу (Panic болдырмау үшін)
	mockAIClient.On("Analyze", mock.Anything, mock.AnythingOfType("ai_client.AnalysisRequest")).
		Return(&ai_client.AnalysisResponse{VisibilityScore: 85.0}, nil)

	// 4. ⚠️ ЖАҢА ҚОСЫЛҒАН ЖОЛ: UpdateResult шақырылымын mock қылу
	// Қате мәтінінде көрсетілгендей: (ctx, 100, "completed", *string)
	// Біз mock.Anything қолданамыз, себебі бізге тек шақырылғаны маңызды
	mockAnalysisRepo.On("UpdateResult", mock.Anything, int64(100), mock.Anything, mock.Anything).Return(nil)

	// Қауіпсіздік үшін UpdateStatus-ты да қалдырамыз (егер ол "processing" деп шақырылса)
	mockAnalysisRepo.On("UpdateStatus", mock.Anything, int64(100), mock.Anything).Return(nil)

	// Сервисті шақыру
	analysisID, err := svc.CreateAnalysis(context.Background(), 1, 1)

	// Нәтижені тексеру
	assert.NoError(t, err)
	assert.Equal(t, int64(100), analysisID)

	// 5. ⚠️ МАҢЫЗДЫ: Goroutine аяқталуы үшін 300 миллисекунд күтеміз
	// Бұл UpdateResult шақырылымының үлгеруіне кепілдік береді
	time.Sleep(300 * time.Millisecond)

	// Барлық mock-тардың шақырылғанын тексеру
	mockCompanyRepo.AssertExpectations(t)
	mockAnalysisRepo.AssertExpectations(t)
	mockAIClient.AssertExpectations(t)
}
