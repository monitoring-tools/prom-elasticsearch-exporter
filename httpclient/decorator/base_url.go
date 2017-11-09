package decorator

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
)

// BaseURLDecorator returns a DecoratorFunc that sends request to all of given hosts
// Decorated client will return response from first responded host
func BaseURLDecorator(baseURLs []string) httpclient.DecoratorFunc {
	parsedURLs := make([]*url.URL, len(baseURLs))
	for i := range baseURLs {
		parsed, err := url.Parse(baseURLs[i])
		if err != nil {
			panic(err)
		}

		parsedURLs[i] = parsed
	}

	return func(c httpclient.IClient) httpclient.IClient {
		return httpclient.ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			type response struct {
				url  *url.URL
				resp *http.Response
				err  error
			}

			respc := make(chan response, len(parsedURLs))

			for _, u := range parsedURLs {
				go func(u *url.URL) {
					reqCopy := new(http.Request)
					*reqCopy = *r

					urlCopy := new(url.URL)
					*urlCopy = *r.URL
					urlCopy.Host = u.Host
					urlCopy.Scheme = u.Scheme

					reqCopy.URL = urlCopy

					res, err := c.Do(reqCopy)
					if err != nil {
						respc <- response{urlCopy, res, err}
						return
					}
					if res.StatusCode != http.StatusOK {
						respc <- response{urlCopy, res, err}
						return
					}

					respc <- response{urlCopy, res, nil}
				}(u)
			}

			errors := []error{}
			for {
				select {
				case r := <-respc:
					if r.err == nil {
						return r.resp, nil
					}
					errors = append(errors, r.err)

					if len(errors) >= len(parsedURLs) {
						return r.resp, fmt.Errorf("Failed to fetch data from all URLs")
					}
				}
			}
		})
	}
}
