package shared

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type ErrorKind uint8

const (
	KindInternal ErrorKind = iota
	KindInvalid
	KindUnauthorized
	KindForbidden
	KindNotFound
	KindConflict
)

var statusByKind = map[ErrorKind]int{
	KindInternal:     http.StatusInternalServerError,
	KindInvalid:      http.StatusBadRequest,
	KindUnauthorized: http.StatusUnauthorized,
	KindForbidden:    http.StatusForbidden,
	KindNotFound:     http.StatusNotFound,
	KindConflict:     http.StatusConflict,
}

type AppError struct {
	kind ErrorKind
	msg  string
}

func newAppError(kind ErrorKind, msg string) *AppError {
	return &AppError{kind: kind, msg: msg}
}

func Invalid(msg string) *AppError      { return newAppError(KindInvalid, msg) }
func Unauthorized(msg string) *AppError { return newAppError(KindUnauthorized, msg) }
func Forbidden(msg string) *AppError    { return newAppError(KindForbidden, msg) }
func NotFound(msg string) *AppError     { return newAppError(KindNotFound, msg) }
func Conflict(msg string) *AppError     { return newAppError(KindConflict, msg) }

func (e *AppError) Error() string { return e.msg }

func (e *AppError) Kind() ErrorKind { return e.kind }

func StatusOf(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return statusByKind[appErr.kind]
	}

	var validationErrs ValidationErrors
	if errors.As(err, &validationErrs) {
		return http.StatusBadRequest
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Code
	}

	return http.StatusInternalServerError
}

func MessageOf(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.msg
	}

	var validationErrs ValidationErrors
	if errors.As(err, &validationErrs) {
		return validationErrs.Error()
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Message
	}

	return internalErrorMessage
}
