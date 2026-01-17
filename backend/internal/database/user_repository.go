package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aslam/backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (d *DB) CreateUser(ctx context.Context, email, password string, role models.UserRole) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()

	query := `
		INSERT INTO users (id, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, role, created_at, updated_at
	`

	user := &models.User{ID: id}
	err = d.conn.QueryRowContext(ctx, query, id, email, string(hashedPassword), string(role)).
		Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	return user, nil
}

func (d *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := d.conn.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DB) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, password, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := d.conn.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DB) ValidatePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func (d *DB) UpdateUser(ctx context.Context, id string, updates *models.UpdateUserRequest) (*models.User, error) {
	query := `
		UPDATE users
		SET email = COALESCE($1, email),
		    role = COALESCE($2, role),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, email, role, created_at, updated_at
	`

	user := &models.User{}
	email := sql.NullString{Valid: updates.Email != ""}
	email.String = updates.Email

	role := sql.NullString{Valid: updates.Role != ""}
	role.String = string(updates.Role)

	err := d.conn.QueryRowContext(ctx, query, email, role, id).
		Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DB) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := d.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (d *DB) ListUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, email, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := d.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
