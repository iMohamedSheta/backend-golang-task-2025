package errors

import (
	"taskgo/pkg/enums"
)

type BadRequestError struct {
	BaseError
}

func NewBadRequestError(publicMsg, privateMsg string, err error) *BadRequestError {
	return &BadRequestError{
		BaseError: newBaseError(enums.ErrCodeBadRequest, publicMsg, privateMsg, err),
	}
}

func AsBadRequestError(err error) (*BadRequestError, bool) {
	se, ok := err.(*BadRequestError)
	return se, ok
}
