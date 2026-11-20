package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"hackathon-absolutecinema/internal/models"
)

// ============ User Repository Mock ============
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, email, passwordHash, name string) (int64, error) {
	args := m.Called(ctx, email, passwordHash, name)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// ============ Company Repository Mock ============
type MockCompanyRepository struct {
	mock.Mock
}

func (m *MockCompanyRepository) Create(ctx context.Context, userID int64, name, website, description string) (int64, error) {
	args := m.Called(ctx, userID, name, website, description)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCompanyRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Company, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Company), args.Error(1)
}

func (m *MockCompanyRepository) GetByID(ctx context.Context, id int64) (*models.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyRepository) Update(ctx context.Context, id int64, name, website, description string) error {
	args := m.Called(ctx, id, name, website, description)
	return args.Error(0)
}

func (m *MockCompanyRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ============ Analysis Repository Mock ============
type MockAnalysisRepository struct {
	mock.Mock
}

func (m *MockAnalysisRepository) Create(ctx context.Context, userID, companyID int64, status string) (int64, error) {
	args := m.Called(ctx, userID, companyID, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAnalysisRepository) GetByID(ctx context.Context, id int64) (*models.Analysis, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Analysis), args.Error(1)
}

func (m *MockAnalysisRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAnalysisRepository) UpdateResult(ctx context.Context, id int64, status string, result *string) error {
	args := m.Called(ctx, id, status, result)
	return args.Error(0)
}
