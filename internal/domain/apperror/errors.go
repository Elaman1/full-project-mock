package apperror

import "errors"

var (
	ExistsEmailErr = errors.New("email already exists")
)
