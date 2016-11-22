/*
Package client is the HTTP(S) client for picel.
*/
package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"time"

	"github.com/henvic/picel/version"
)

var (
	client = &http.Client{}

	// UserAgent for the picel middleware
	UserAgent = "picel/" + version.Version + " (+https://github.com/henvic/picel)"

	// ErrBackend is a generic error returned when the server fails to fulfill the request
	ErrBackend = errors.New("Backend server failed to fulfill the request")
)

// Download a given URL
type Download struct {
	URL           string
	Filename      string
	file          *os.File
	request       *http.Request
	timeout       *time.Duration
	cancelTimeout *context.CancelFunc
	context       context.Context
}

// Timeout for the request
func (d *Download) Timeout(timeout time.Duration) {
	d.timeout = &timeout
}

// Cancel the download
func (d *Download) Cancel() {
	if d.cancelTimeout != nil {
		(*d.cancelTimeout)()
	}
}

// Load the download
func (d *Download) Load() (err error) {
	if err = d.createFile(); err != nil {
		return err
	}

	defer d.file.Close()

	if err = d.setupRequest(); err != nil {
		return err
	}

	d.request.Header.Set("User-Agent", UserAgent)
	return d.do()
}

func (d *Download) createFile() (err error) {
	d.file, err = os.Create(d.Filename)
	return err
}

func (d *Download) setupRequestTimeout() {
	if d.timeout != nil && *d.timeout != 0*time.Second {
		var c context.CancelFunc
		d.context, c = context.WithTimeout(d.context, *d.timeout)
		d.cancelTimeout = &c
		d.request = d.request.WithContext(d.context)
	}
}

func (d *Download) setupRequest() (err error) {
	d.request, err = http.NewRequest("GET", d.URL, nil)

	if err != nil {
		return err
	}

	d.context = context.Background()
	d.setupRequestTimeout()
	d.request = d.request.WithContext(d.context)
	return nil
}

func (d *Download) do() (err error) {
	var resp *http.Response
	resp, err = client.Do(d.request)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		_, err = io.Copy(d.file, resp.Body)
		return err
	case http.StatusNotFound:
		return http.ErrMissingFile
	}

	return ErrBackend
}
