package nibo

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	userPasswordCost int = bcrypt.DefaultCost
)

var (
	errShortPassword error = errors.New("password must be > 3 characters")
)

// UserRepository is a database abstraction for performing CRUD operations on users.
type UserRepository interface {
	Create(context.Context, CreateUser) (*User, error)
	GetByID(context.Context, string) (*User, error)
	GetByPassword(context.Context, string) (*User, error)
	Update(context.Context, string, UpdateUser) (*User, error)
	Delete(context.Context, string) error
}

// UpdateUser represents a request to update an existing user.
type UpdateUser struct {
	// DisplayName is the new display name to set.
	DisplayName string

	// Password is the new password to set.
	Password string
}

// CreateUser represents a request to create a new user.
type CreateUser struct {
	// Name is the user's registered name. See User.Name.
	Name string

	// DisplayName is the user's alternate name. See User.DisplayName.
	DisplayName string

	// EmailAddress is the user's email address. See User.EmailAddress.
	EmailAddress string

	// Password is the desired password to register the new user with.
	Password string
}

// User represents a user account.
type User struct {
	// ID is the user's unique ID.
	ID string

	// Name is the user's registered name. Names are unique across all users.
	Name string

	// DisplayName is the user's alternate name. This name, if defined, is prefered for presentation
	// over Name.
	DisplayName string

	// EmailAddress is the user's registered email address. Email addresses are unique across all
	// users.
	EmailAddress string
}

// GeneratePassword generates an app-specific hash given a plaintext password.
func GeneratePassword(plaintext string) ([]byte, error) {
	if len(plaintext) < 3 {
		return nil, errShortPassword
	}

	h, err := bcrypt.GenerateFromPassword([]byte(plaintext), userPasswordCost)
	if err != nil {
		return nil, err
	}
	return h, nil
}
