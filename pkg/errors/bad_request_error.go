package errors

import (
	"fmt"
	"taskgo/pkg/enums"
)

type BadRequestError struct {
	ErrorCode      string
	PublicMessage  string
	PrivateMessage string
	Err            error
}

func NewBadRequestError(publicMessage string, privateMessage string, err error) *BadRequestError {
	if publicMessage == "" {
		publicMessage = "Bad Request"
	}
	if privateMessage == "" {
		privateMessage = "Bad Request"
	}

	return &BadRequestError{
		ErrorCode:      string(enums.ErrCodeBadRequest),
		PublicMessage:  publicMessage,
		PrivateMessage: privateMessage,
		Err:            err,
	}
}

func (e *BadRequestError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("BadRequestError: %s", e.PrivateMessage)
	}
	return fmt.Sprintf("BadRequestError: %s - %v", e.PrivateMessage, e.Err)
}

func (e *BadRequestError) PublicError() string {
	return e.PublicMessage
}

func (e *BadRequestError) Unwrap() error {
	return e.Err
}

func AsBadRequestError(err error) (*BadRequestError, bool) {
	se, ok := err.(*BadRequestError)
	return se, ok
}
