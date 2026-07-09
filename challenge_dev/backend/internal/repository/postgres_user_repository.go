package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tianpl92/retojikko_sebaspaniagua/challenge_dev/backend/internal/domain"
)

// PostgresUserRepository is a PostgreSQL-backed implementation of UserRepository.
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgresUserRepository.
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

// FindByID retrieves a user by document number (primary key).
func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, gender, email, phone_number, password,
		       status, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	row := r.pool.QueryRow(ctx, query, id)
	user, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}

// FindByEmail retrieves a user by email.
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, gender, email, phone_number, password,
		       status, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	row := r.pool.QueryRow(ctx, query, email)
	user, err := scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}

// Create inserts a new user into the database.
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, gender, email, phone_number, password,
		                   status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Gender,
		user.Email, user.Phone, user.Password, user.Status,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		if isPGDuplicate(err) {
			return ErrDuplicateEmail
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// ExistsByDocument checks if a user with the given document number exists.
func (r *PostgresUserRepository) ExistsByDocument(ctx context.Context, document string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, document).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("exists by document: %w", err)
	}
	return exists, nil
}

// ExistsByEmail checks if a user with the given email exists.
func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("exists by email: %w", err)
	}
	return exists, nil
}

// Update updates an existing user's data.
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, gender = $4, email = $5,
		    phone_number = $6, password = $7, status = $8, updated_at = $9
		WHERE id = $1 AND deleted_at IS NULL
	`
	ct, err := r.pool.Exec(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Gender,
		user.Email, user.Phone, user.Password, user.Status,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// scanUser scans a pgx.Row into a domain.User.
func scanUser(row pgx.Row) (*domain.User, error) {
	user := &domain.User{}
	err := row.Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Gender,
		&user.Email, &user.Phone, &user.Password,
		&user.Status, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// isPGDuplicate checks if a PostgreSQL error is a unique constraint violation.
func isPGDuplicate(err error) bool {
	const pgUniqueViolation = "23505"
	var pgErr *pgconn.PgError
	if as, ok := err.(*pgconn.PgError); ok {
		pgErr = as
	} else {
		return false
	}
	return pgErr.Code == pgUniqueViolation
}
