package errors

import (
	"fmt"
	"taskgo/pkg/enums"
)

type UnAuthorizedError struct {
	ErrorCode      string
	PublicMessage  string
	PrivateMessage string
	Err            error
}

func NewUnAuthorizedError(publicMsg, privateMsg string, err error) *UnAuthorizedError {
	if publicMsg == "" {
		publicMsg = "Unauthorized Action"
	}
	if privateMsg == "" {
		privateMsg = publicMsg
	}
	return &UnAuthorizedError{
		ErrorCode:      string(enums.ErrCodeUnauthorized),
		PublicMessage:  publicMsg,
		PrivateMessage: privateMsg,
		Err:            err,
	}
}

func (e *UnAuthorizedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("UnauthorizedError: %s - %v", e.PrivateMessage, e.Err)
	}
	return fmt.Sprintf("UnauthorizedError: %s", e.PrivateMessage)
}

func (e *UnAuthorizedError) PublicError() string {
	return e.PublicMessage
}

func (e *UnAuthorizedError) Unwrap() error {
	return e.Err
}

func AsUnAuthorizedError(err error) (*UnAuthorizedError, bool) {
	authorizeErr, ok := err.(*UnAuthorizedError)
	return authorizeErr, ok
}
