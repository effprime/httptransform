package httptransform

import (
	"net/http"

	"github.com/gorilla/mux"
)

type TransformOptions struct {
	RequestTransformer         RequestTransformer
	SuccessResponseTransformer ResponseTransformer
	ErrorResponseTransformer   ResponseTransformer
	ErrorHandler               ErrorHandler
}

func (o *TransformOptions) getErrorHandler() ErrorHandler {
	var errorHandler ErrorHandler
	errorHandler = &defaultErrorHandler{}
	if o.ErrorHandler != nil {
		errorHandler = o.ErrorHandler
	}
	return errorHandler
}

func Transform(opts *TransformOptions) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req, err := opts.RequestTransformer.TransformRequest(req)
			if err != nil {
				opts.getErrorHandler().Handle(err, w, req)
				return
			}
			next.ServeHTTP(NewResponseTransformWriter(&ResponseTransformerWriterOptions{
				SuccessTransformer: opts.SuccessResponseTransformer,
				ErrorTransformer:   opts.ErrorResponseTransformer,
				Writer:             w,
			}), req)
		})
	}
}
