package user

import (
	"errors"
	"full-project-mock/pkg/validator"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	if r.Username == "" || r.Password == "" || r.Email == "" {
		return errors.New("username, email or password is empty")
	}

	if !validator.IsValidEmail(r.Email) {
		return errors.New("invalid email")
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if r.Email == "" || r.Password == "" {
		return errors.New("email or password is empty")
	}

	if !validator.IsValidEmail(r.Email) {
		return errors.New("invalid email")
	}

	return nil
}
