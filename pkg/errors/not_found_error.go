package errors

import (
	"taskgo/pkg/enums"
)

type NotFoundError struct {
	BaseError
}

func NewNotFoundError(publicMsg, privateMsg string, err error) *NotFoundError {
	return &NotFoundError{
		BaseError: newBaseError(enums.ErrCodeNotFound, publicMsg, privateMsg, err),
	}
}

func AsNotFoundError(err error) (*NotFoundError, bool) {
	nfe, ok := err.(*NotFoundError)
	return nfe, ok
}
