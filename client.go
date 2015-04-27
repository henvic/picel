package main

import (
	"errors"
	"io"
	"net/http"
	"os"
)

var (
	ErrBackend = errors.New("Backend server failed to fulfill the request")
)

func Load(url string, filename string) (size int64, err error) {
	file, err := os.Create(filename)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return io.Copy(file, resp.Body)
	case http.StatusNotFound:
		return 0, http.ErrMissingFile
	default:
		return 0, ErrBackend
	}
}
