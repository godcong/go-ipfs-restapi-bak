package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/godcong/go-ipfs-files"
	"io"
	"os"
	"path"
	"path/filepath"
)

// AddRet ...
type AddRet struct {
	Hash string
	Name string
	Size string
}

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

// AddFile ...
func (a *API) AddFile(pathname string) (*AddRet, error) {
	stat, err := os.Lstat(pathname)
	if err != nil {
		return nil, err
	}

	sf, err := files.NewSerialFile(pathname, false, stat)
	if err != nil {
		return nil, err
	}

	_, file := filepath.Split(pathname)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry(path.Base(file), sf)})
	fileReader := files.NewMultiFileReader(slf, true)

	var out AddRet
	req := a.Request("add")
	req.Body = fileReader

	e := req.Exec(context.Background(), &out)
	return &out, e
}

//Add ...
func (a *API) Add(r io.Reader, options ...AddOpts) (*AddRet, error) {
	fr := files.NewReaderFile(r)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	var out AddRet
	req := a.Request("add")

	for _, option := range options {
		_ = option(req)
	}
	req.Body = fileReader
	e := req.Exec(context.Background(), &out)
	return &out, e
}

// AddLink ...
func (a *API) AddLink(target string) (*AddRet, error) {
	link := files.NewLinkFile(target, nil)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", link)})
	fileReader := files.NewMultiFileReader(slf, true)

	var out AddRet
	req := a.Request("add")
	req.Body = fileReader

	e := req.Exec(context.Background(), &out)
	return &out, e
}

// AddDir adds a directory recursively with all of the files under it
func (a *API) AddDir(dir string) ([]*AddRet, error) {
	stat, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}

	sf, err := files.NewSerialFile(dir, false, stat)
	if err != nil {
		return nil, err
	}

	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry(path.Base(dir), sf)})
	reader := files.NewMultiFileReader(slf, true)

	req := a.Request("add").Option("recursive", true)
	req.Body = reader
	responder, err := req.POST(context.Background())
	if err != nil {
		return nil, err
	}

	defer responder.Close()

	if responder.Error != nil {
		return nil, responder.Error
	}

	dec := json.NewDecoder(responder.Output)
	var final []*AddRet

	for {
		var out AddRet
		err = dec.Decode(&out)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		final = append(final, &out)
	}

	if final == nil {
		return nil, errors.New("no results received")
	}

	return final, nil
}
