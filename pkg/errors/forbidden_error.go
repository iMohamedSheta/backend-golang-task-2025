package errors

import (
	"taskgo/pkg/enums"
)

type ForbiddenError struct {
	BaseError
}

func NewForbiddenError(publicMsg, privateMsg string, err error) *ForbiddenError {
	return &ForbiddenError{
		BaseError: newBaseError(enums.ErrCodeForbidden, publicMsg, privateMsg, err),
	}
}

func AsForbiddenError(err error) (*ForbiddenError, bool) {
	forbiddenErr, ok := err.(*ForbiddenError)
	return forbiddenErr, ok
}
