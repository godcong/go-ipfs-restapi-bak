package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// API ...
type API interface {
}

// api ...
type api struct {
	url     string
	client  *http.Client
	opts    map[string]string
	headers map[string]string
	body    io.Reader
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
	requester := buildRequester(a.url, command, args...)
	requester.Opts = a.opts
	requester.Body = a.body
	requester.Headers = a.headers
	requester.Client = a.client
	return requester.Do(context.Background())
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
