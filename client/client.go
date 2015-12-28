/*
Package client is the HTTP(S) client for picel.
*/
package client

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/henvic/picel/version"
)

var (
	client     = &http.Client{}
	UserAgent  = "picel/" + version.Version + " (+https://github.com/henvic/picel)"
	ErrBackend = errors.New("Backend server failed to fulfill the request")
)

func Load(url string, filename string) (size int64, err error) {
	file, err := os.Create(filename)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return 0, err
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return io.Copy(file, resp.Body)
	case http.StatusNotFound:
		return 0, http.ErrMissingFile
	}

	return 0, ErrBackend
}
