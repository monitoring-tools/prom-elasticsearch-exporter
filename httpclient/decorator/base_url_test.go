package decorator

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
	. "gopkg.in/check.v1"
)

var baseURLs = []string{"http://host-1", "http://host-2"}
var baseURLsPassed = make([]bool, len(baseURLs))

func (s *TestSuite) TestBaseURLDecorator(c *C) {
	baseURLCheckerDecorator := func(c httpclient.IClient) httpclient.IClient {
		return httpclient.ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			for i := range baseURLs {
				parsed, _ := url.Parse(baseURLs[i])

				if r.URL.Host == parsed.Host {
					baseURLsPassed[i] = true
				}
			}

			return nil, errors.New("")
		})
	}
	r, _ := http.NewRequest("GET", "/", nil)
	httpClient := httpclient.Decorate(s.dummyClient, baseURLCheckerDecorator, BaseURLDecorator(baseURLs))
	httpClient.Do(r)

	for i := range baseURLsPassed {
		c.Assert(baseURLsPassed[i], Equals, true)
	}
}
