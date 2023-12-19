package models

import (
	"errors"
)

var (
	//Given snippet not found.
	ErrNoRecord = errors.New("models: no matching record found")
	//User tries to log in with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	//User tries to log in with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)