package api

import (
	"fmt"
	"net/http"
	"strings"
)

// API ...
type API interface {
}

// api ...
type api struct {
	url    string
	client *http.Client
}

// New ...
func New(url string) API {
	c := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}

	return NewWithClient(url, c)
}

// NewWithClient ...
func NewWithClient(url string, c *http.Client) API {
	var api api
	api.url = url
	api.client = c
	// We don't support redirects.
	api.client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return fmt.Errorf("unexpected redirect")
	}
	return &api
}

// GET ...
func (*api) GET() {

}

// POST ...
func (*api) POST() {

}

// Request ...
func (a *api) Request(command string, args ...string) (*Responder, error) {
	return buildRequester(a.url, command, args...).Do(a.client)
}

func buildRequester(url, command string, args ...string) *Requester {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	opts := map[string]string{
		"encoding":        "json",
		"stream-channels": "true",
	}

	return &Requester{
		ApiBase: url + "/api/v0",
		Command: command,
		Args:    args,
		Opts:    opts,
		Headers: make(map[string]string),
	}
}
