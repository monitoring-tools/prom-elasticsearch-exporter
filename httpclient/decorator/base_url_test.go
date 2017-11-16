package decorator

import (
	"errors"
	"net/http"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestBaseURLDecorator(c *C) {
	baseURLCheckerDecorator := func(client httpclient.IClient) httpclient.IClient {
		return httpclient.ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			c.Assert(r.URL.Host, Equals, "my-host")
			c.Assert(r.URL.Scheme, Equals, "https")

			return nil, errors.New("")
		})
	}
	r, _ := http.NewRequest("GET", "/", nil)
	httpClient := httpclient.Decorate(s.dummyClient, baseURLCheckerDecorator, BaseURLDecorator("https://my-host"))
	httpClient.Do(r)
}
