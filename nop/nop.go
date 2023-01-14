package nop

import "net/http"

// NopRequestTransformer is a request transformer that applies no transformation (no-op).
type NopRequestTransformer struct{}

// TransformRequest returns the provided request with no transformation.
func (t *NopRequestTransformer) TransformRequest(req *http.Request) (*http.Request, error) {
	return req, nil
}

// NopResponseTransformer is a response transformer that applies no transformation (no-op).
type NopResponseTransformer struct{}

// TransformResponse returns the provided response with no transformation.
func (t *NopResponseTransformer) TransformResponse(b []byte) ([]byte, error) {
	return b, nil
}
