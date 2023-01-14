package httptransform

import (
	"net/http"

	"github.com/effprime/httptransform/nop"
	"github.com/gorilla/mux"
)

type TransformOptions struct {
	// The transformer applied to all incoming requests.
	// If not provided, a nop transformer is used.
	RequestTransformer RequestTransformer

	// The transformer applied to all success (200-level) responses.
	// If no error response transformer is provided, non-200 responses are handled as well.
	// If not provided, a nop transformer is used.
	SuccessResponseTransformer ResponseTransformer

	// The transformer applied to all non-200 level responses.
	// If not provided, the success transformer is used for error transformations.
	ErrorResponseTransformer ResponseTransformer

	// Interface for handling transformation errors.
	// If not provided, the default error handler is used.
	ErrorHandler ErrorHandler
}

// Transform returns a middleware handler that transforms incoming requests and outgoing responses.
// For any errors in transforming, an error repsonse is generated using either the provided error handler
// or the default one if not provided.
func Transform(opts *TransformOptions) mux.MiddlewareFunc {
	opts.fillDefaults()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req, err := opts.RequestTransformer.TransformRequest(req)
			if err != nil {
				opts.ErrorHandler.HandleTransformError(err, w)
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

func (o *TransformOptions) fillDefaults() {
	if o.RequestTransformer == nil {
		o.RequestTransformer = &nop.NopRequestTransformer{}
	}
	if o.SuccessResponseTransformer == nil {
		o.SuccessResponseTransformer = &nop.NopResponseTransformer{}
	}
	if o.ErrorResponseTransformer == nil {
		o.ErrorResponseTransformer = o.SuccessResponseTransformer
	}
	if o.ErrorHandler == nil {
		o.ErrorHandler = &defaultErrorHandler{}
	}
}
