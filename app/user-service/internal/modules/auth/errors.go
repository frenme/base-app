package auth

import "errors"

var (
	ErrorInvalidCredentials  = errors.New("invalid username or password")
	ErrorInvalidRefreshToken = errors.New("invalid refresh token")
	ErrorInvalidPassword     = errors.New("password length must be between 4 and 16 characters")
	ErrorUserAlreadyExists   = errors.New("user already exists")
)
