package errors

import (
	"taskgo/pkg/enums"
)

type ServerError struct {
	BaseError
}

func NewServerError(publicMsg, privateMsg string, err error) *ServerError {
	return &ServerError{
		BaseError: newBaseError(enums.ErrCodeInternalError, publicMsg, privateMsg, err),
	}
}

func AsServerError(err error) (*ServerError, bool) {
	se, ok := err.(*ServerError)
	return se, ok
}
