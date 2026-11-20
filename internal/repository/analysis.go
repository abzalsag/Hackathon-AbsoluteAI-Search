package repository

import (
	"context"
	"hackathon-absolutecinema/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type analysisRepository struct {
	pool *pgxpool.Pool
}

func NewAnalysisRepository(pool *pgxpool.Pool) AnalysisRepository {
	return &analysisRepository{pool: pool}
}

func (r *analysisRepository) Create(ctx context.Context, userID, companyID int64, status string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, "INSERT INTO analyses (user_id, company_id, status) VALUES ($1, $2, $3) RETURNING id", userID, companyID, status).Scan(&id)
	return id, err
}

func (r *analysisRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.pool.Exec(ctx, "UPDATE analyses SET status = $1 WHERE id = $2", status, id)
	return err
}

func (r *analysisRepository) UpdateResult(ctx context.Context, id int64, status string, result *string) error {
	_, err := r.pool.Exec(ctx, "UPDATE analyses SET status = $1, result = $2 WHERE id = $3", status, result, id)
	return err
}

func (r *analysisRepository) GetByID(ctx context.Context, id int64) (*models.Analysis, error) {
	var a models.Analysis
	err := r.pool.QueryRow(ctx, "SELECT id, user_id, company_id, status, result, created_at FROM analyses WHERE id = $1", id).
		Scan(&a.ID, &a.UserID, &a.CompanyID, &a.Status, &a.Result, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
