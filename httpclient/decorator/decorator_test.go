package decorator

import (
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

var (
	_ = Suite(&TestSuite{})
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
	dummyClient *dummy
}

type dummy struct{}

func (d *dummy) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{Request: r}, nil
}
