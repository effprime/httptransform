package httptransform

import "net/http"

type ErrorHandler interface {
	Handle(error, http.ResponseWriter, *http.Request)
}

type defaultErrorHandler struct{}

func (h *defaultErrorHandler) Handle(err error, w http.ResponseWriter, req *http.Request) {
	tErr, ok := err.(*TransformError)
	if !ok {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	http.Error(w, tErr.Error(), tErr.StatusCode())
}

type TransformError struct {
	message string
	code    int
}

func (e *TransformError) Error() string {
	return e.message
}

func (e *TransformError) StatusCode() int {
	return e.code
}
