package api

import (
	"context"
	"github.com/ipfs/go-ipfs-files"
	"io"
)

// AddOpts ...
type AddOpts = func(requester *Requester) error

// OnlyHash ...
func OnlyHash(enabled bool) AddOpts {
	return func(rb *Requester) error {
		rb.Option("only-hash", enabled)
		return nil
	}
}

// Pin ...
func Pin(enabled bool) AddOpts {
	return func(rb *Requester) error {
		rb.Option("pin", enabled)
		return nil
	}
}

// Progress ...
func Progress(enabled bool) AddOpts {
	return func(rb *Requester) error {
		rb.Option("progress", enabled)
		return nil
	}
}

// RawLeaves ...
func RawLeaves(enabled bool) AddOpts {
	return func(rb *Requester) error {
		rb.Option("raw-leaves", enabled)
		return nil
	}
}

//Add ...
func (a *api) Add(r io.Reader, options ...AddOpts) (map[string]string, error) {
	fr := files.NewReaderFile(r)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	var out map[string]string
	req := a.Request("add")

	for _, option := range options {
		_ = option(req)
	}
	req.Body = fileReader
	e := req.Exec(context.Background(), &out)
	return out, e
}
