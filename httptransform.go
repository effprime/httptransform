package httptransform

import (
	"net/http"
)

type RequestTransformer interface {
	TransformRequest(*http.Request) (*http.Request, error)
}

type ResponseTransformer interface {
	TransformResponse([]byte) ([]byte, error)
}
