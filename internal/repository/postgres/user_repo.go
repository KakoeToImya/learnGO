package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/KakoeToImya/go-ws-chat/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, username, password, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW()) 
		RETURNING id`

	err := r.pool.QueryRow(ctx, query,
		user.Email, user.Username, user.Password,
	).Scan(&user.ID)

	if err != nil {
		// Проверяем на дубликат
		if isUniqueViolation(err) {
			return fmt.Errorf("user with this email or username already exists")
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at 
		FROM users WHERE id = $1`

	return r.scanUser(r.pool.QueryRow(ctx, query, id))
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at 
		FROM users WHERE email = $1`

	return r.scanUser(r.pool.QueryRow(ctx, query, email))
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at 
		FROM users WHERE username = $1`

	return r.scanUser(r.pool.QueryRow(ctx, query, username))
}

func (r *UserRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at 
		FROM users ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username,
			&user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return users, nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, username = $2, password = $3, updated_at = NOW()
		WHERE id = $4`

	_, err := r.pool.Exec(ctx, query,
		user.Email, user.Username, user.Password, user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// Вспомогательные методы

func (r *UserRepo) scanUser(row pgx.Row) (*domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID, &user.Email, &user.Username,
		&user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &user, nil
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "23505")
}
