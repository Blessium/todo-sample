package errors

import (
	"time"
)

type ErrorType uint8

const (
	ErrUnknown ErrorType = iota
    ErrInternal
	ErrValidation
	ErrNotExist
)

func (e ErrorType) String() string {
	switch e {
	case ErrInternal:
		return "internal error"
	case ErrValidation:
		return "schema is not valid"
	case ErrNotExist:
		return "element doesn't exist"
	default:
		return "unknown error"
	}
}

type HttpErrorContent struct {
	Status  uint   `json:"status"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

func NewHttpErrorResponse(status uint, typ string, message string, path string) HttpErrorResponse {
    content := HttpErrorContent {
        Status: status,
        Type: typ,
        Message: message,
        Path: path,
    }

    return HttpErrorResponse {
        Error : content,
    }
}


type HttpErrorResponse struct {
	Error HttpErrorContent `json:"error"`
}

type Error struct {
	From      string
	Type      ErrorType
	Timestamp time.Time
	Message   string
	Details   error
}

func NewError(from string, errorType ErrorType, message string, details error) error {
	return Error{
		From:      from,
		Type:      errorType,
		Timestamp: time.Now(),
		Message:   message,
		Details:   details,
	}
}

func (e Error) Error() string {
    return e.Message
}

func IsType(err error, t ErrorType) bool {
    e, ok := err.(Error)
    if !ok {
        return false
    }
    return e.Type == t
}
