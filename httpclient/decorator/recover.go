package decorator

import (
	"fmt"
	"net/http"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
)

// RecoverDecorator returns a DecoratorFunc that recovers panic and convert it to error
func RecoverDecorator() httpclient.DecoratorFunc {
	return func(c httpclient.Client) httpclient.Client {
		return httpclient.ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("recovered panic: %v", r)
				}
			}()
			return c.Do(r)
		})
	}
}
