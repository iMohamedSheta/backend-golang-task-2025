package errors

import (
	"fmt"
	"taskgo/pkg/enums"
)

type NotFoundError struct {
	ErrorCode      string
	PublicMessage  string
	PrivateMessage string
	Err            error
}

func NewNotFoundError(publicMessage string, privateMessage string, err error) *NotFoundError {
	if publicMessage == "" {
		publicMessage = "Not found"
	}
	if privateMessage == "" {
		privateMessage = "Not found"
	}

	return &NotFoundError{
		ErrorCode:      string(enums.ErrCodeNotFound),
		PublicMessage:  publicMessage,
		PrivateMessage: privateMessage,
		Err:            err,
	}
}

func (e *NotFoundError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("NotFoundError: %s", e.PrivateMessage)
	}
	return fmt.Sprintf("NotFoundError: %s - %v", e.PrivateMessage, e.Err)
}

func (e *NotFoundError) PublicError() string {
	return e.PublicMessage
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

func AsNotFoundError(err error) (*NotFoundError, bool) {
	nfe, ok := err.(*NotFoundError)
	return nfe, ok
}
