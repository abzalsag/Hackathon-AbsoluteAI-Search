package repository

import (
	"context"
	"hackathon-absolutecinema/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type companyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) CompanyRepository {
	return &companyRepository{pool: pool}
}

func (r *companyRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Company, error) {
	var companies []models.Company
	rows, err := r.pool.Query(ctx, "SELECT id, user_id, name, website, description, created_at FROM companies WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Website, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

func (r *companyRepository) Create(ctx context.Context, userID int64, name, website, description string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, "INSERT INTO companies (user_id, name, website, description) VALUES ($1, $2, $3, $4) RETURNING id", userID, name, website, description).Scan(&id)
	return id, err
}

func (r *companyRepository) GetByID(ctx context.Context, id int64) (*models.Company, error) {
	var c models.Company
	err := r.pool.QueryRow(ctx, "SELECT id, user_id, name, website, description, created_at FROM companies WHERE id = $1", id).
		Scan(&c.ID, &c.UserID, &c.Name, &c.Website, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *companyRepository) Update(ctx context.Context, id int64, name, website, description string) error {
	_, err := r.pool.Exec(ctx, "UPDATE companies SET name = $1, website = $2, description = $3 WHERE id = $4", name, website, description, id)
	return err
}

func (r *companyRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM companies WHERE id = $1", id)
	return err
}
