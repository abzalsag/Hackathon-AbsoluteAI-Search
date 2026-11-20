package repository

import (
	"context"
	"hackathon-absolutecinema/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, email, passwordHash, name string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, "INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id", email, passwordHash, name).Scan(&id)
	return id, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.pool.QueryRow(ctx, "SELECT id, email, password_hash, name, created_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
