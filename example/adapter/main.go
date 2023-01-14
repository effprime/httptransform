package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/effprime/httptransform"
	"github.com/gorilla/mux"
)

// This example shows a basic request adapter solution.
//
// The incoming request (ClientRequest) identifies a person and their address,
// where the address is defined by multiple sub-fields.
// The handler that is assumed un-changable accepts an APIRequest,
// where the address is defined as a single string.
// The transformer will change the request to an APIRequest.

type ClientRequest struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
	Street string `json:"street"`
	City   string `json:"city"`
}

type APIRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type DB struct{}

func (db *DB) SaveRequest(name, address string) error {
	// some implementation of saving the request
	return nil
}

func HandleAPIRequest(w http.ResponseWriter, req *http.Request) {
	db := DB{}
	body := APIRequest{}
	json.NewDecoder(req.Body).Decode(&body)
	db.SaveRequest(body.Name, body.Address)
}

type ClientRequestTransformer struct{}

func (t *ClientRequestTransformer) TransformRequest(req *http.Request) (*http.Request, error) {
	clientRequest := ClientRequest{}
	err := json.NewDecoder(req.Body).Decode(&clientRequest)
	if err != nil {
		return req, httptransform.NewTransformError(err.Error(), http.StatusBadRequest).WithExternal("could not read client request")
	}

	apiRequest := APIRequest{
		Name:    clientRequest.Name,
		Address: fmt.Sprintf("%v %v, %v", clientRequest.Number, clientRequest.Street, clientRequest.City),
	}

	output, err := json.Marshal(apiRequest)
	if err != nil {
		return req, httptransform.NewTransformError(err.Error(), http.StatusInternalServerError).WithExternal("could not serialize transformed response")
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(output))
	req.ContentLength = int64(len(output))

	return req, nil
}

func main() {
	r := mux.NewRouter()

	handler := httptransform.Transform(&httptransform.TransformOptions{
		RequestTransformer: &ClientRequestTransformer{},
	})(http.HandlerFunc(HandleAPIRequest))

	r.Handle("/save", handler).Methods(http.MethodPost)
	http.ListenAndServe(":8000", r)
}
