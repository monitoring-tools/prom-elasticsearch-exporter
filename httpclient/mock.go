package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client implementation of IClient for testing purpose
type Client struct {
	promises []Promise
}

// NewClient creates new client
func NewClient() *Client {
	return &Client{}
}

// Do comparing request with all constrains and returns first matching
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	requestBody := readerToString(r.Body)
	for _, promise := range c.promises {
		// Request copy is required because body can be read only once
		requestCopy := new(http.Request)
		*requestCopy = *r
		requestCopy.Body = closer{strings.NewReader(requestBody)}
		if response, ok := promise.Check(requestCopy); ok {
			return response, nil
		}
	}
	return nil, fmt.Errorf("No suitable request found")
}

// Get is a helper method for using PromiseBuilder
func (c *Client) Get(path string) *PromiseBuilder {
	return NewPromiseBuilder(c).Get(path)
}

// Post is a helper method for using PromiseBuilder
func (c *Client) Post(path string) *PromiseBuilder {
	return NewPromiseBuilder(c).Post(path)
}

// Request is a helper method that returns empty PromiseBuilder
// It should be used if you just want to mock response
// and don't need to check request method, path and/or body
func (c *Client) Request() *PromiseBuilder {
	return NewPromiseBuilder(c)
}

func readerToString(r io.Reader) string {
	if r == nil {
		return ""
	}
	buf := &bytes.Buffer{}
	buf.ReadFrom(r)

	return buf.String()
}
