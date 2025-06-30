package errors

import (
	"taskgo/pkg/enums"
)

type UnAuthorizedError struct {
	BaseError
}

func NewUnAuthorizedError(publicMsg, privateMsg string, err error) *UnAuthorizedError {
	return &UnAuthorizedError{
		BaseError: newBaseError(enums.ErrCodeUnauthorized, publicMsg, privateMsg, err),
	}
}

func AsUnAuthorizedError(err error) (*UnAuthorizedError, bool) {
	authorizeErr, ok := err.(*UnAuthorizedError)
	return authorizeErr, ok
}
