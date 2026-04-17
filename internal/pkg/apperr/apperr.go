package apperr

import "fmt"

type AppError struct {
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"`
	Data       any    `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func New(message string, httpStatus int) *AppError {
	return &AppError{Message: message, HTTPStatus: httpStatus}
}

func (e *AppError) Wrap(err error) *AppError {
	return &AppError{Message: e.Message, HTTPStatus: e.HTTPStatus, Err: err, Data: e.Data}
}

func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{Message: msg, HTTPStatus: e.HTTPStatus, Err: e.Err, Data: e.Data}
}

func (e *AppError) WithData(data any) *AppError {
	return &AppError{Message: e.Message, HTTPStatus: e.HTTPStatus, Err: e.Err, Data: data}
}
