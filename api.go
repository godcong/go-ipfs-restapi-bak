package api

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"strings"
)

// API ...
type API struct {
	url     string
	client  *http.Client
	opts    map[string]string
	headers map[string]string
	body    io.Reader
}

// New ...
func New(url string) *API {
	c := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}

	return NewWithClient(url, c)
}

// NewWithClient ...
func NewWithClient(url string, c *http.Client) *API {
	var api API
	api.url = url
	api.client = c
	// We don't support redirects.
	api.client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return fmt.Errorf("unexpected redirect")
	}
	return &api
}

// Request ...
func (a *API) Request(command string, args ...string) *Requester {
	requester := buildRequester(a.url, command, args...)
	requester.Opts = a.opts
	requester.Body = a.body
	requester.Headers = a.headers
	requester.Client = a.client
	return requester
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
		APIBase: url + "/api/v0",
		Command: command,
		Args:    args,
		Opts:    opts,
		Headers: make(map[string]string),
	}
}

// AddDirList ...
func (a *API) AddDirList(path string) (*ListObject, error) {
	rets, e := a.AddDir(path)
	if e != nil {
		return &ListObject{}, e
	}
	size := len(rets)
	if size == 0 {
		return &ListObject{}, xerrors.New("add not result")
	}
	return a.List("/ipfs/" + rets[len(rets)-1].Hash)
}

// List entries at the given path
func (a *API) List(path string) (*ListObject, error) {
	var out struct {
		Objects []*ListObject
	}
	err := a.Request("ls", path).Exec(context.Background(), &out)
	if err != nil {
		return nil, err
	}
	if len(out.Objects) != 1 {
		return nil, errors.New("bad response from server")
	}
	return out.Objects[0], nil
}

// ListLink ...
type ListLink struct {
	Hash string
	Name string
	Size uint64
	Type int
}

// ListObject ...
type ListObject struct {
	Links []*ListLink
	ListLink
}

// Pin the given path
func (a *API) Pin(path string) error {
	return a.Request("pin/add", path).
		Option("recursive", true).
		Exec(context.Background(), nil)
}

// Unpin the given path
func (a *API) Unpin(path string) error {
	return a.Request("pin/rm", path).
		Option("recursive", true).
		Exec(context.Background(), nil)
}
