package service

import (
	"context"

	"hackathon-absolutecinema/internal/models"
	"hackathon-absolutecinema/internal/repository"
)

type CompanyService struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyService(companyRepo repository.CompanyRepository) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
	}
}

func (s *CompanyService) GetCompanies(ctx context.Context, userID int64) ([]models.Company, error) {
	return s.companyRepo.GetByUserID(ctx, userID)
}

func (s *CompanyService) CreateCompany(ctx context.Context, userID int64, name, website, description string) (int64, error) {
	return s.companyRepo.Create(ctx, userID, name, website, description)
}

func (s *CompanyService) GetCompany(ctx context.Context, id int64) (*models.Company, error) {
	return s.companyRepo.GetByID(ctx, id)
}

func (s *CompanyService) UpdateCompany(ctx context.Context, id int64, name, website, description string) error {
	return s.companyRepo.Update(ctx, id, name, website, description)
}

func (s *CompanyService) DeleteCompany(ctx context.Context, id int64) error {
	return s.companyRepo.Delete(ctx, id)
}
