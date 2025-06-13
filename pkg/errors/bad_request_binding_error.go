package errors

import (
	"fmt"
	"taskgo/pkg/enums"
)

type BadRequestBindingError struct {
	ErrorCode      string
	PublicMessage  string
	PrivateMessage string
	Err            error
}

func NewBadRequestBindingError(publicMessage string, privateMessage string, err error) *BadRequestBindingError {
	if publicMessage == "" {
		publicMessage = "Bad Request"
	}
	if privateMessage == "" {
		privateMessage = "Bad Request"
	}

	return &BadRequestBindingError{
		ErrorCode:      string(enums.ErrCodeBadRequest),
		PublicMessage:  publicMessage,
		PrivateMessage: privateMessage,
		Err:            err,
	}
}

func (e *BadRequestBindingError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("BadRequestBindingError: %s", e.PrivateMessage)
	}
	return fmt.Sprintf("BadRequestBindingError: %s - %v", e.PrivateMessage, e.Err)
}

func (e *BadRequestBindingError) PublicError() string {
	return e.PublicMessage
}

func (e *BadRequestBindingError) Unwrap() error {
	return e.Err
}

func AsBadRequestBindingError(err error) (*BadRequestBindingError, bool) {
	se, ok := err.(*BadRequestBindingError)
	return se, ok
}
