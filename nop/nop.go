package nop

import "net/http"

type NopRequestTransformer struct{}

func (t *NopRequestTransformer) TransformRequest(req *http.Request) (*http.Request, error) {
	return req, nil
}

type NopResponseTransformer struct{}

func (t *NopResponseTransformer) TransformRequest(b []byte) ([]byte, error) {
	return b, nil
}
