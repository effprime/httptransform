# httptransform

```bash
go get github.com/effprime/httptransform
```

`httptransform` is a middleware library for transforming incoming HTTP requests and outgoing HTTP responses.

First, build a basic transformer:

```go
// Define a struct that satifies request transformer interface
type ClientRequestTransformer struct{}

// Implement request transformation
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
```

Then attach it to the middleware:

```go
func main() {
    r := mux.NewRouter()

    // Attach middleware to some handler that accepts transformed request
    handler := httptransform.Transform(&httptransform.TransformOptions{
        RequestTransformer: &ClientRequestTransformer{},
    })(http.HandlerFunc(HandleAPIRequest))

    r.Handle("/save", handler).Methods(http.MethodPost)
    http.ListenAndServe(":8000", r)
}
```


For more complete examples, see `examples/`


