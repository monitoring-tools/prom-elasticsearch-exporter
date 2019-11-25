package decorator

import (
	"net/http"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestPanicDecorator(c *C) {
	panicTriggerDecorator := func(c httpclient.Client) httpclient.Client {
		return httpclient.ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			panic("oops")
		})
	}
	r, _ := http.NewRequest("GET", "/", nil)
	httpClient := httpclient.Decorate(s.dummyClient, panicTriggerDecorator, RecoverDecorator())
	res, err := httpClient.Do(r)

	c.Assert(res, IsNil)
	c.Assert(err, ErrorMatches, "*oops")
}
