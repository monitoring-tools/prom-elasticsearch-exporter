package httpclient

import (
	"testing"

	"bytes"
	"net/http"

	"io"

	. "gopkg.in/check.v1"
)

type TestSuite struct{}

var (
	_ = Suite(&TestSuite{})
)

type dummyClient struct{}

func (c *dummyClient) Do(r *http.Request) (*http.Response, error) {
	return nil, nil
}

func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) writeHeaders(headers http.Header) string {
	buf := &bytes.Buffer{}
	headers.Write(buf)

	return buf.String()
}

func (s *TestSuite) builder() *PromiseBuilder {
	return NewPromiseBuilder(NewClient())
}

func (s *TestSuite) newRequest(method, path string, body io.Reader) *http.Request {
	r, _ := http.NewRequest(method, path, body)

	return r
}
