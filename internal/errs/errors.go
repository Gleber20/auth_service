package errs

import "errors"

var (
	ErrNotfound                    = errors.New("not found")
	ErrUserNotFound                = errors.New("user not found")
	ErrInvalidUserID               = errors.New("invalid user id")
	ErrInvalidRequestBody          = errors.New("invalid request body")
	ErrInvalidFieldValue           = errors.New("invalid field value")
	ErrUsernameAlreadyExists       = errors.New("username already exists")
	ErrIncorrectUsernameOrPassword = errors.New("incorrect username or password")
	ErrInvalidToken                = errors.New("invalid token")
	ErrSomethingWentWrong          = errors.New("something went wrong")
)
