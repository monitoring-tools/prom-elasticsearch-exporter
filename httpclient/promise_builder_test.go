package httpclient

import (
	"bytes"
	"net/http"
	"strings"

	. "gopkg.in/check.v1"
)

type TestRequestBuilderSuite struct {
	TestSuite
}

var _ = Suite(&TestRequestBuilderSuite{})

func (s *TestRequestBuilderSuite) TestNewConstraintBuilder(c *C) {
	builder := s.builder()

	c.Assert(len(builder.checkers), Equals, 0)
}

func (s *TestRequestBuilderSuite) TestWithMethod(c *C) {
	builder := s.builder().WithMethod("PUT")

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("PUT", "/", nil),
		s.newRequest("PUT", "/path/", nil),
		s.newRequest("PUT", "/path/", strings.NewReader("something")),
	}, true)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("GET", "/", nil),
		s.newRequest("POST", "/", nil),
	}, false)
}

func (s *TestRequestBuilderSuite) TestGet(c *C) {
	builder := s.builder().Get("/some/path")

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("GET", "/some/path", nil),
		s.newRequest("GET", "/some/path", strings.NewReader("something")),
	}, true)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("POST", "/some/path", nil),
		s.newRequest("GET", "/some/path/", nil),
	}, false)
}

func (s *TestRequestBuilderSuite) TestPost(c *C) {
	builder := s.builder().Post("/some/path")

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("POST", "/some/path", nil),
		s.newRequest("POST", "/some/path", strings.NewReader("something")),
	}, true)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("GET", "/some/path", nil),
		s.newRequest("GET", "/some/path/", nil),
	}, false)
}

func (s *TestRequestBuilderSuite) TestWithPath(c *C) {
	builder := s.builder().WithPath("/some/path")

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("POST", "/some/path", nil),
		s.newRequest("GET", "/some/path", nil),
	}, true)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("GET", "/some/path/", nil),
		s.newRequest("GET", "/", nil),
	}, false)
}

func (s *TestRequestBuilderSuite) TestWithHeader(c *C) {
	builder := s.builder().WithHeader("X-Foo", "123")

	r1 := s.newRequest("POST", "/some/path", nil)
	r1.Header.Add("X-Foo", "123")

	r2 := s.newRequest("POST", "/some/path", nil)
	r2.Header.Add("X-Foo", "1234")

	r3 := s.newRequest("POST", "/some/path", nil)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		r1,
	}, true)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		r2,
		r3,
	}, false)
}

func (s *TestRequestBuilderSuite) TestWithHost(c *C) {
	builder := s.builder().WithHost("some.host")

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("POST", "http://some.host/path", nil),
		s.newRequest("GET", "http://some.host/path", nil),
		s.newRequest("GET", "//some.host/some/path/", nil),
	}, true)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("GET", "/", nil),
		s.newRequest("GET", "some.host", nil),
		s.newRequest("GET", "some.host/some/path/", nil),
	}, false)
}

func (s *TestRequestBuilderSuite) TestWithBody(c *C) {
	builder := s.builder().WithBodyReader(strings.NewReader("some body"))

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("POST", "/some/path", strings.NewReader("some body")),
		s.newRequest("PUT", "/some/path", strings.NewReader("some body")),
	}, true)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("POST", "/some/path", nil),
		s.newRequest("PUT", "/some/path", strings.NewReader("wrong body")),
	}, false)
}

func (s *TestRequestBuilderSuite) TestBody(c *C) {
	builder := s.builder().WithBody("some body")

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("POST", "/some/path", strings.NewReader("some body")),
		s.newRequest("PUT", "/some/path", strings.NewReader("some body")),
	}, true)

	s.assertCheckers(c, builder.checkers, []*http.Request{
		s.newRequest("POST", "/some/path", nil),
		s.newRequest("PUT", "/some/path", strings.NewReader("wrong body")),
	}, false)
}

func (s *TestRequestBuilderSuite) TestReturns(c *C) {
	client := NewClient()
	builder := NewPromiseBuilder(client)
	builder.WillReturn(200, "some response")

	c.Assert(len(client.promises), Equals, 1)
	c.Assert(client.promises[0].checkers, DeepEquals, builder.checkers)

	buf := &bytes.Buffer{}
	buf.ReadFrom(client.promises[0].response.Body)

	c.Assert(client.promises[0].response.StatusCode, Equals, 200)
	c.Assert(buf.String(), Equals, "some response")
}

func (s *TestRequestBuilderSuite) TestReturnsResponse(c *C) {
	response := &http.Response{StatusCode: 200, Body: &closer{strings.NewReader("some response")}}

	client := NewClient()
	builder := NewPromiseBuilder(client)
	builder.WillReturnResponse(response)

	c.Assert(len(client.promises), Equals, 1)
	c.Assert(client.promises[0].checkers, DeepEquals, builder.checkers)
	c.Assert(client.promises[0].response, Equals, response)
}

func (s *TestRequestBuilderSuite) TestCloser(c *C) {
	// Fake test. Just for 100% coverage
	closer := &closer{}
	c.Assert(closer.Close(), IsNil)
}

func (s *TestRequestBuilderSuite) assertCheckers(c *C, checkers []CheckerFunc, requests []*http.Request, valid bool) {
	rule := &Promise{checkers: checkers}
	for _, req := range requests {
		_, ok := rule.Check(req)
		c.Assert(ok, Equals, valid)
	}
}
