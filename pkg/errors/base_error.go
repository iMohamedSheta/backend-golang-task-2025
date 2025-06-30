package errors

import (
	"fmt"
	"taskgo/pkg/enums"
)

type BaseError struct {
	ErrorCode      enums.ErrorCode `json:"code"`
	PublicMessage  string          `json:"message"`
	PrivateMessage string          `json:"-"`
	Err            error           `json:"-"`
}

func (e *BaseError) Error() string {
	return formatErrorMessage(e.ErrorCode, e.PrivateMessage, e.Err)
}

func (e *BaseError) PublicError() string {
	return formatErrorMessage(e.ErrorCode, e.PublicMessage, nil)
}

func (e *BaseError) Unwrap() error {
	return e.Err
}

// Reuse this across all error types
func newBaseError(code enums.ErrorCode, publicMsg, privateMsg string, err error) BaseError {
	if publicMsg == "" {
		publicMsg = code.Message()
	}
	if privateMsg == "" {
		privateMsg = code.Message()
	}

	return BaseError{
		ErrorCode:      code,
		PublicMessage:  publicMsg,
		PrivateMessage: privateMsg,
		Err:            err,
	}
}

// Helper func for formatting the error messages
func formatErrorMessage(code enums.ErrorCode, message string, err error) string {
	baseMessage := message
	if baseMessage == "" {
		baseMessage = code.Message()
	}

	if err != nil {
		return fmt.Sprintf("%s: %s - %v", code, baseMessage, err)
	}

	return fmt.Sprintf("%s: %s", code, baseMessage)
}
