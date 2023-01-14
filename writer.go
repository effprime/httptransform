package httptransform

import (
	"net/http"

	"github.com/effprime/httptransform/nop"
)

type responseTransformWriter struct {
	successTansformer ResponseTransformer
	errorTransformer  ResponseTransformer
	writer            http.ResponseWriter
	code              int
	errorHandler      ErrorHandler
}

type ResponseTransformerWriterOptions struct {
	// The transformer applied to all success (200-level) responses.
	// If no error response transformer is provided, non-200 responses are handled as well.
	// If not provided, a nop transformer is used.
	SuccessTransformer ResponseTransformer

	// The transformer applied to all non-200 level responses.
	// If not provided, the success transformer is used for error transformations.
	ErrorTransformer ResponseTransformer

	// The underlying HTTPResponseWriter.
	Writer http.ResponseWriter

	// Handler for response transformation errors.
	ErrorHandler ErrorHandler
}

func (o *ResponseTransformerWriterOptions) fillDefaults() {
	if o.SuccessTransformer == nil {
		o.SuccessTransformer = &nop.NopResponseTransformer{}
	}
	if o.ErrorTransformer == nil {
		o.ErrorTransformer = &nop.NopResponseTransformer{}
	}
	if o.ErrorHandler == nil {
		o.ErrorHandler = &defaultErrorHandler{}
	}
}

// NewResponseTransformWriter returns a new response transforming http.ResponseWriter.
func NewResponseTransformWriter(opts *ResponseTransformerWriterOptions) *responseTransformWriter {
	return &responseTransformWriter{
		successTansformer: opts.SuccessTransformer,
		errorTransformer:  opts.ErrorTransformer,
		writer:            opts.Writer,
		errorHandler:      opts.ErrorHandler,
	}
}

// Write attempts to transform outgoing responses and then writing them to the underlying writer.
// If the response is 200-level, responses are transformed via the SuccessTransformer.
// Otherwise, they are transformed via the ErrorTransformer.
func (w *responseTransformWriter) Write(inc []byte) (int, error) {
	var err error
	output := inc
	if w.success() {
		output, err = w.successTansformer.TransformResponse(inc)
		if err != nil {
			return w.errorHandler.HandleTransformError(err, w.writer)
		}
	} else {
		output, err = w.errorTransformer.TransformResponse(inc)
		if err != nil {
			return w.errorHandler.HandleTransformError(err, w.writer)
		}
	}
	return w.writer.Write(output)
}

func (w *responseTransformWriter) success() bool {
	return w.code/100 == 2
}

// WriteHeader stores the status code internally and then calls WriteHeader on the underlying writer.
func (w *responseTransformWriter) WriteHeader(code int) {
	w.code = code
	w.writer.WriteHeader(code)
}

// Header returns the underlying writer's header.
func (w *responseTransformWriter) Header() http.Header {
	return w.writer.Header()
}
