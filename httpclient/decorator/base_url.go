package decorator

import (
	"net/http"
	"net/url"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
)

// BaseURLDecorator returns a DecoratorFunc that modifies request
// and sets given baseURL as request host
func BaseURLDecorator(baseURL string) httpclient.DecoratorFunc {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	return func(c httpclient.IClient) httpclient.IClient {
		return httpclient.ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			reqCopy := new(http.Request)
			*reqCopy = *r

			urlCopy := new(url.URL)
			*urlCopy = *r.URL
			urlCopy.Host = parsed.Host
			urlCopy.Scheme = parsed.Scheme

			reqCopy.URL = urlCopy

			return c.Do(reqCopy)
		})
	}
}
