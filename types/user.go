package types

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// DeleteUserInput represents the input for deleting a user
type DeleteUserInput struct {
	ID int `json:"id"`
}

// ReadUserInput represents the input for reading a user
type ReadUserInput struct {
	ID int `json:"id"`
}

// CreateUserOutput represents the output of creating a user
type CreateUserOutput struct {
	User User `json:"user"`
}

// ReadUserOutput represents the output of reading a user
type ReadUserOutput struct {
	User User `json:"user"`
}

// UpdateUserOutput represents the output of updating a user
type UpdateUserOutput struct {
	User User `json:"user"`
}

// DeleteUserOutput represents the output of deleting a user
type DeleteUserOutput struct {
	Success bool `json:"success"`
}
