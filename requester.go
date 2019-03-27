package api

import (
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-ipfs-files"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Requester ...
type Requester struct {
	ApiBase string
	Command string
	Args    []string
	Opts    map[string]string
	Body    io.Reader
	Headers map[string]string
}

// Result ...
type Responder struct {
	Decoder map[string]interface{}
	Command string
	Message string
	Code    int
}

// Do ...
func (r *Requester) Do(client *http.Client) (responder *Responder, e error) {
	req, e := http.NewRequest("POST", r.URL(), r.Body)
	if e != nil {
		return nil, e
	}

	// Add any headers that were supplied via the RequestBuilder.
	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}

	if fr, ok := r.Body.(*files.MultiFileReader); ok {
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+fr.Boundary())
		req.Header.Set("Content-Disposition", "form-data; name=\"files\"")
	}

	resp, e := client.Do(req)
	if e != nil {
		return nil, e
	}

	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, ";")
	contentType = parts[0]
	responder = &Responder{}
	if resp.StatusCode >= http.StatusBadRequest {
		responder.Command = r.Command
		switch {
		case resp.StatusCode == http.StatusNotFound:
			responder.Message = "command not found"
		case contentType == "text/plain":
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response (%d) read error: %s\n", resp.StatusCode, err)
			}
			responder.Message = string(out)
		case contentType == "application/json":
			//if e = json.NewDecoder(resp.Body).Decode(e); e != nil {
			//	_, _ = fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response (%d) unmarshall error: %s\n", resp.StatusCode, e)
			//}
			//body, err := ioutil.ReadAll(resp.Body)
			dec := json.NewDecoder(resp.Body)
			//var fileInfo []map[string]string
			//out := map[string]string{}
			for {
				e = dec.Decode(&responder)
				if e != nil {
					if e == io.EOF {
						break
					}
					return nil, e
				}
				//m = out
				//fileInfo = append(fileInfo, out)
			}

		default:
			_, _ = fmt.Fprintf(os.Stderr, "ipfs-shell: warning! unhandled response (%d) encoding: %s", resp.StatusCode, contentType)
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "ipfs-shell: response (%d) read error: %s\n", resp.StatusCode, err)
			}
			responder.Message = fmt.Sprintf("unknown ipfs-shell error encoding: %q - %q", contentType, out)
		}
		//responder.Error = e
		//responder.Output = nil

		// drain body and close
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return &Responder{}, nil
}

// URL ...
func (r *Requester) URL() string {
	values := make(url.Values)
	for _, arg := range r.Args {
		values.Add("arg", arg)
	}
	for k, v := range r.Opts {
		values.Add(k, v)
	}

	return fmt.Sprintf("%s/%s?%s", r.ApiBase, r.Command, values.Encode())
}
