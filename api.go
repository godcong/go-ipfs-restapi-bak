package api

import (
	"fmt"
	"net/http"
)

// Api ...
type Api struct {
	url    string
	client *http.Client
}

// New ...
func New(url string) *Api {
	c := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}

	return NewWithClient(url, c)
}

// NewWithClient ...
func NewWithClient(url string, c *http.Client) *Api {
	var api Api
	api.url = url
	api.client = c
	// We don't support redirects.
	api.client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return fmt.Errorf("unexpected redirect")
	}
	return &api
}

// GET ...
func (*Api) GET() {

}

// POST ...
func (*Api) POST() {

}

// Request ...
func (*Api) Request() {
	buildRequester()
}

func buildRequester() {

}
