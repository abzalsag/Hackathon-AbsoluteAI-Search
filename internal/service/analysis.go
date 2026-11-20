package service

import (
	"context"
	"encoding/json"
	"fmt"

	"hackathon-absolutecinema/internal/ai_client"
	"hackathon-absolutecinema/internal/models"
	"hackathon-absolutecinema/internal/repository"
)

type AnalysisService struct {
	analysisRepo repository.AnalysisRepository // Интерфейс
	companyRepo  repository.CompanyRepository  // Интерфейс
	aiClient     ai_client.Client              // Интерфейс
}

func NewAnalysisService(
	analysisRepo repository.AnalysisRepository,
	companyRepo repository.CompanyRepository,
	aiClient ai_client.Client,
) *AnalysisService {
	return &AnalysisService{
		analysisRepo: analysisRepo,
		companyRepo:  companyRepo,
		aiClient:     aiClient,
	}
}

func (s *AnalysisService) CreateAnalysis(ctx context.Context, userID, companyID int64) (int64, error) {
	company, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return 0, fmt.Errorf("компания табылмады: %w", err)
	}
	analysisID, err := s.analysisRepo.Create(ctx, userID, companyID, "pending")
	if err != nil {
		return 0, err
	}
	go s.processAnalysisAsync(analysisID, company)

	return analysisID, nil
}

func (s *AnalysisService) processAnalysisAsync(analysisID int64, company *models.Company) {
	ctx := context.Background()
	_ = s.analysisRepo.UpdateStatus(ctx, analysisID, "processing")
	req := ai_client.AnalysisRequest{
		CompanyID:   company.ID,
		Name:        company.Name,
		Website:     company.Website,
		Description: company.Description,
	}

	result, err := s.aiClient.Analyze(ctx, req)
	if err != nil {
		_ = s.analysisRepo.UpdateStatus(ctx, analysisID, "failed")
		return
	}
	resultJSON, _ := json.Marshal(result)
	resultStr := string(resultJSON)

	_ = s.analysisRepo.UpdateResult(ctx, analysisID, "completed", &resultStr)
}

func (s *AnalysisService) GetAnalysis(ctx context.Context, id int64) (*models.Analysis, error) {
	return s.analysisRepo.GetByID(ctx, id)
}
