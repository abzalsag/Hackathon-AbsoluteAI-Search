package service

import (
	"context"
	"hackathon-absolutecinema/internal/models"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, email, password, name string) (int64, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type CompanyServiceInterface interface {
	GetCompanies(ctx context.Context, userID int64) ([]models.Company, error)
	CreateCompany(ctx context.Context, userID int64, name, website, description string) (int64, error)
	GetCompany(ctx context.Context, id int64) (*models.Company, error)
	UpdateCompany(ctx context.Context, id int64, name, website, description string) error
	DeleteCompany(ctx context.Context, id int64) error
}

type AnalysisServiceInterface interface {
	CreateAnalysis(ctx context.Context, userID, companyID int64) (int64, error)
	GetAnalysis(ctx context.Context, id int64) (*models.Analysis, error)
}
