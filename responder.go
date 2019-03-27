package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Responder ...
type Responder struct {
	Output io.ReadCloser
	Error  error
}

// Close ...
func (r *Responder) Close() (e error) {
	if r.Output != nil {
		// always drain output (response body)
		_, e := io.Copy(ioutil.Discard, r.Output)
		if e != nil {
			return e
		}
		e = r.Output.Close()
		if e != nil {
			return e
		}
	}
	return nil
}

// Decode ...
func (r *Responder) Decode(dec interface{}) (e error) {
	defer r.Close()
	if r.Error != nil {
		return r.Error
	}

	return json.NewDecoder(r.Output).Decode(dec)
}
