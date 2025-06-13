package errors

import (
	"fmt"
	"taskgo/pkg/enums"
)

type ServerError struct {
	ErrorCode      string
	PublicMessage  string
	PrivateMessage string
	Err            error
}

func NewServerError(publicMessage string, privateMessage string, err error) *ServerError {
	if publicMessage == "" {
		publicMessage = "Internal Server Error"
	}
	if privateMessage == "" {
		privateMessage = "Internal Server Error"
	}

	return &ServerError{
		ErrorCode:      string(enums.ErrCodeInternalError),
		PublicMessage:  publicMessage,
		PrivateMessage: privateMessage,
		Err:            err,
	}
}

func (e *ServerError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("ServerError: %s", e.PrivateMessage)
	}
	return fmt.Sprintf("ServerError: %s - %v", e.PrivateMessage, e.Err)
}

func (e *ServerError) PublicError() string {
	return e.PublicMessage
}

func (e *ServerError) Unwrap() error {
	return e.Err
}

func AsServerError(err error) (*ServerError, bool) {
	se, ok := err.(*ServerError)
	return se, ok
}
