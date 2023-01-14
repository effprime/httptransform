package httptransform

import (
	"fmt"
	"net/http"
)

// Interface for handling transformation errors.
// HandleTransformError should return the number of bytes written to the writer argument,
// Ideally through directly returning calls to w.Write()
type ErrorHandler interface {
	HandleTransformError(error, http.ResponseWriter) (int, error)
}

type defaultErrorHandler struct{}

// HandleTransformError is the default transformation error handler.
// It tries to convert all incoming errors to the defined TransformError,
// which includes an HTTP status code and external message.
// If this conversion fails, the returned response is HTTP 425 with the error message.
func (h *defaultErrorHandler) HandleTransformError(err error, w http.ResponseWriter) (int, error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	tErr, ok := err.(*transformError)
	if !ok {
		w.WriteHeader(http.StatusFailedDependency)
		return w.Write([]byte(fmt.Sprintf("transform error: %v", err.Error())))
	}

	w.WriteHeader(tErr.StatusCode())
	return w.Write([]byte(tErr.ExternalError()))
}

type transformError struct {
	message  string
	external string
	code     int
}

// NewTransformError returns a Transform Error.
// The external error message for HTTP responses is set to the provided message.
func NewTransformError(message string, code int) *transformError {
	return &transformError{
		message:  message,
		external: message,
		code:     code,
	}
}

// Error returns the internal error message.
// This mesgsage is appropriate for logging or other internal cases.
func (e *transformError) Error() string {
	return e.message
}

// ExternalError returns the external error message.
// This message is by default the same as the Error() message,
// But can be set to a more appropriate message for HTTP responses or other external cases.
func (e *transformError) ExternalError() string {
	return e.external
}

// StatusCode returns the HTTP status code associated with the error.
func (e *transformError) StatusCode() int {
	return e.code
}

// WithExternal adds an external error message to the error.
// This message is by default the same as the Error() message,
// But can be set to a more appropriate message for HTTP responses or other external cases.
func (e *transformError) WithExternal(message string) *transformError {
	return &transformError{
		message:  e.message,
		external: message,
		code:     e.code,
	}
}
