package httptransform

import (
	"net/http"
)

// Interface for transforming incoming HTTP requests.
type RequestTransformer interface {
	TransformRequest(*http.Request) (*http.Request, error)
}

// Interface for transforming outgoing HTTP responses.
type ResponseTransformer interface {
	TransformResponse([]byte) ([]byte, error)
}
