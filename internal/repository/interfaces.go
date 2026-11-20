package repository

import (
	"context"
	"hackathon-absolutecinema/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash, name string) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type CompanyRepository interface {
	GetByUserID(ctx context.Context, userID int64) ([]models.Company, error)
	Create(ctx context.Context, userID int64, name, website, description string) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Company, error)
	Update(ctx context.Context, id int64, name, website, description string) error
	Delete(ctx context.Context, id int64) error
}

type AnalysisRepository interface {
	Create(ctx context.Context, userID, companyID int64, status string) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateResult(ctx context.Context, id int64, status string, result *string) error
	GetByID(ctx context.Context, id int64) (*models.Analysis, error)
}
