package httptransform

import "net/http"

type ResponseTransformWriter struct {
	successTansformer ResponseTransformer
	errorTransformer  ResponseTransformer
	writer            http.ResponseWriter
	code              int
}

type ResponseTransformerWriterOptions struct {
	SuccessTransformer ResponseTransformer
	ErrorTransformer   ResponseTransformer
	Writer             http.ResponseWriter
}

func NewResponseTransformWriter(opts *ResponseTransformerWriterOptions) *ResponseTransformWriter {
	return &ResponseTransformWriter{
		successTansformer: opts.SuccessTransformer,
		errorTransformer:  opts.ErrorTransformer,
		writer:            opts.Writer,
	}
}

func (w *ResponseTransformWriter) Write(inc []byte) (int, error) {
	var err error
	output := inc
	if w.Success() {
		output, err = w.successTansformer.TransformResponse(inc)
		if err != nil {
			return 0, err
		}
	} else {
		output, err = w.errorTransformer.TransformResponse(inc)
		if err != nil {
			return 0, err
		}
	}
	return w.writer.Write(output)
}

func (w *ResponseTransformWriter) Success() bool {
	return w.code/100 == 2
}

func (w *ResponseTransformWriter) WriteHeader(code int) {
	w.code = code
	w.writer.WriteHeader(code)
}

func (w *ResponseTransformWriter) Header() http.Header {
	return w.writer.Header()
}
