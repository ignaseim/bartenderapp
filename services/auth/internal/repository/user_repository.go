package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/ignaseim/bartenderapp/services/pkg/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE user_id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// List retrieves all users, optionally filtered by role
func (r *UserRepository) List(role string) ([]models.User, error) {
	var query string
	var rows *sql.Rows
	var err error

	if role != "" {
		query = `
			SELECT user_id, username, email, password_hash, role, created_at, updated_at
			FROM users
			WHERE role = $1
			ORDER BY username
		`
		rows, err = r.db.Query(query, role)
	} else {
		query = `
			SELECT user_id, username, email, password_hash, role, created_at, updated_at
			FROM users
			ORDER BY username
		`
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Create adds a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	// Check if user exists
	_, err := r.GetByID(user.ID)
	if err != nil {
		return err
	}

	// If password is being updated
	if user.PasswordHash != "" {
		query := `
			UPDATE users
			SET username = $1, email = $2, password_hash = $3, role = $4, updated_at = NOW()
			WHERE user_id = $5
			RETURNING updated_at
		`

		err = r.db.QueryRow(
			query,
			user.Username,
			user.Email,
			user.PasswordHash,
			user.Role,
			user.ID,
		).Scan(&user.UpdatedAt)
	} else {
		// If password is not being updated
		query := `
			UPDATE users
			SET username = $1, email = $2, role = $3, updated_at = NOW()
			WHERE user_id = $4
			RETURNING updated_at
		`

		err = r.db.QueryRow(
			query,
			user.Username,
			user.Email,
			user.Role,
			user.ID,
		).Scan(&user.UpdatedAt)
	}

	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}

	return nil
}

// Delete removes a user
func (r *UserRepository) Delete(id int) error {
	query := `
		DELETE FROM users
		WHERE user_id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
} 