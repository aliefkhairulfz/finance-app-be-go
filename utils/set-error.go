package utils

import "errors"

var (
	ErrorNoRowsFound    = errors.New("no rows found")
	InternalServerError = errors.New("internal server error")

	ErrorNoUserFound = errors.New("no user found")

	ErrorEmailConflict = errors.New("email already exists")
	ErrorEmailNotFound = errors.New("email not found")

	ErrorAccountConflict = errors.New("account already exists")
	ErrorAccountNotFound = errors.New("account not found")

	ErrorAccountWrongPassword = errors.New("wrong password account")
	ErrorSessionNotFound      = errors.New("session not found")
)
