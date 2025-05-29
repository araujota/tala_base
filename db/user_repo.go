package db

import (
	"database/sql"
	"fmt"

	"tala_base/types"
)

// CreateUser creates a new user in the database.
// This function is called by the user_create lambda to persist user data.
// It returns the created user with its ID and timestamps.
func CreateUser(db *sql.DB, input types.CreateUserInput) (*types.User, error) {
	var user types.User
	err := db.QueryRow(
		`INSERT INTO users (email, name) 
		VALUES ($1, $2) 
		RETURNING id, email, name, created_at, updated_at`,
		input.Email, input.Name,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID.
// This function is called by the user_read lambda to fetch user details.
// It returns a user if found, or an error if not found or on database error.
func GetUserByID(db *sql.DB, id int) (*types.User, error) {
	var user types.User
	err := db.QueryRow(
		`SELECT id, email, name, created_at, updated_at 
		FROM users 
		WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// ListUsers retrieves all users from the database.
// This function is called by the user_list lambda to fetch all users.
// It returns a slice of users, or an error if the database query fails.
func ListUsers(db *sql.DB) ([]*types.User, error) {
	rows, err := db.Query(
		`SELECT id, email, name, created_at, updated_at 
		FROM users 
		ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*types.User
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}

// UpdateUser updates an existing user's information.
// This function is called by the user_update lambda to modify user data.
// It returns the updated user with new timestamps.
func UpdateUser(db *sql.DB, id int, input types.UpdateUserInput) (*types.User, error) {
	var user types.User
	err := db.QueryRow(
		`UPDATE users 
		SET email = $1, name = $2 
		WHERE id = $3 
		RETURNING id, email, name, created_at, updated_at`,
		input.Email, input.Name, id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &user, nil
}

// DeleteUser removes a user from the database.
// This function is called by the user_delete lambda to remove a user.
// It returns an error if the user is not found or if the deletion fails.
func DeleteUser(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found: %d", id)
	}
	return nil
}
