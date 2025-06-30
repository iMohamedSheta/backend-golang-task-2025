package errors

import (
	"taskgo/pkg/enums"
)

type BadRequestBindingError struct {
	BaseError
}

func NewBadRequestBindingError(publicMsg, privateMsg string, err error) *BadRequestBindingError {
	return &BadRequestBindingError{
		BaseError: newBaseError(enums.ErrCodeBadRequest, publicMsg, privateMsg, err),
	}
}

func AsBadRequestBindingError(err error) (*BadRequestBindingError, bool) {
	se, ok := err.(*BadRequestBindingError)
	return se, ok
}
