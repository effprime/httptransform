package main

import (
	"encoding/json"
	"net/http"

	"github.com/effprime/httptransform"
	"github.com/gorilla/mux"
)

// This example shows how to convert a response to JSON.
//
// The outgoing response is a simple "Hello world" string,
// so the transformer will wrap this text in JSON with field "message"

func HelloWorldHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello world"))
}

type JSONResponse struct {
	Message string `json:"message"`
}

type JSONResponseTransformer struct{}

func (t *JSONResponseTransformer) TransformResponse(inc []byte) ([]byte, error) {
	output, err := json.Marshal(&JSONResponse{
		Message: string(inc),
	})
	if err != nil {
		return nil, httptransform.NewTransformError("could not serialize JSON response", http.StatusInternalServerError)
	}
	return output, nil
}

func main() {
	r := mux.NewRouter()

	handler := httptransform.Transform(&httptransform.TransformOptions{
		SuccessResponseTransformer: &JSONResponseTransformer{},
	})(http.HandlerFunc(HelloWorldHandler))

	r.Handle("/hello", handler)
	http.ListenAndServe(":8000", r)
}
