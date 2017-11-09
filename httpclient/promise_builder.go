package httpclient

import (
	"io"
	"net/http"
	"strings"
)

// CheckerFunc checks if request fits conditions
type CheckerFunc func(r *http.Request) bool

// Promise is a combination of some checker functions and response.
type Promise struct {
	checkers []CheckerFunc
	response *http.Response
}

// PromiseBuilder helper for creating rules with fluent interface
type PromiseBuilder struct {
	client   *Client
	checkers []CheckerFunc
}

type closer struct {
	io.Reader
}

// Check check requests passes all checkers
func (p *Promise) Check(r *http.Request) (*http.Response, bool) {
	var valid bool

	requestBody := readerToString(r.Body)
	// Request copy is required because body can be read only once
	requestCopy := new(http.Request)
	*requestCopy = *r
	requestCopy.Body = closer{strings.NewReader(requestBody)}

	valid = true

	for _, checker := range p.checkers {
		if !checker(requestCopy) {
			valid = false
			break
		}
	}

	if !valid {
		return nil, false
	}

	return p.response, true
}

func (c closer) Close() error {
	return nil
}

// NewPromiseBuilder constructor for ruleBuilder.
// all created constrains are registered in client
func NewPromiseBuilder(client *Client) *PromiseBuilder {
	return &PromiseBuilder{
		client:   client,
		checkers: make([]CheckerFunc, 0),
	}
}

// Get checks that HTTP Method equals to GET
func (b *PromiseBuilder) Get(path string) *PromiseBuilder {
	return b.WithMethod("GET").WithPath(path)
}

// Post checks that HTTP Method equals to POST
func (b *PromiseBuilder) Post(path string) *PromiseBuilder {
	return b.WithMethod("POST").WithPath(path)
}

// WithMethod checks that HTTP Method equals to specified method
func (b *PromiseBuilder) WithMethod(method string) *PromiseBuilder {
	return b.WithChecker(func(req *http.Request) bool {
		return req.Method == method
	})
}

// WithPath checks that path equals to specified
func (b *PromiseBuilder) WithPath(path string) *PromiseBuilder {
	return b.WithChecker(func(req *http.Request) bool {
		return req.URL.RequestURI() == path
	})
}

// WithHeader checks that request contains specified header with specified value
func (b *PromiseBuilder) WithHeader(key, value string) *PromiseBuilder {
	return b.WithChecker(func(req *http.Request) bool {
		return req.Header.Get(key) == value
	})
}

// WithHost checks that request host equals to specified
func (b *PromiseBuilder) WithHost(host string) *PromiseBuilder {
	return b.WithChecker(func(req *http.Request) bool {
		return req.Host == host
	})
}

// WithBodyReader checks that body equals to specified
func (b *PromiseBuilder) WithBodyReader(body io.Reader) *PromiseBuilder {
	expectedBody := readerToString(body)
	return b.WithChecker(func(req *http.Request) bool {
		requestBody := readerToString(req.Body)

		return requestBody == expectedBody
	})
}

// WithChecker adds custorm checker
func (b *PromiseBuilder) WithChecker(checker CheckerFunc) *PromiseBuilder {
	b.checkers = append(b.checkers, checker)

	return b
}

// WithBody checks that body equals to specified
func (b *PromiseBuilder) WithBody(body string) *PromiseBuilder {
	return b.WithBodyReader(strings.NewReader(body))
}

// WillReturn registers builded rule in client
func (b *PromiseBuilder) WillReturn(code int, body string) {
	b.WillReturnResponse(&http.Response{
		StatusCode: code,
		Body:       &closer{strings.NewReader(body)},
	})
}

// WillReturnResponse registers builded rule in client
func (b *PromiseBuilder) WillReturnResponse(response *http.Response) {
	b.client.promises = append(b.client.promises, Promise{
		checkers: b.checkers,
		response: response,
	})
}
