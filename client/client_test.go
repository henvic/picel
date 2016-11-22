package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type LoadProvider struct {
	word string
}

func TestLoadWithInvalidFilename(t *testing.T) {
	var download = &Download{
		URL: "0/foo.png",
	}
	if err := download.Load(); err == nil {
		t.Errorf("Load(0/foo.png, ) should fail")
	}
}

func TestLoadFromInvalidSchema(t *testing.T) {
	file, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
	defer os.Remove(file.Name())

	if tmpFileErr != nil {
		panic(tmpFileErr)
	}

	var download = &Download{
		URL:      "0/foo.png",
		Filename: file.Name(),
	}
	if err := download.Load(); err == nil {
		t.Errorf("Load(0/foo.png, %v) should fail", file.Name())
	}
}

func TestBackendFailure(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	file, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
	defer os.Remove(file.Name())

	if tmpFileErr != nil {
		panic(tmpFileErr)
	}

	var download = &Download{
		URL:      ts.URL,
		Filename: file.Name(),
	}

	if err := download.Load(); err != ErrBackend {
		t.Errorf("Load() should fail with %v, got %v instead", ErrBackend, err)
	}
}

func TestNotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	file, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
	defer os.Remove(file.Name())

	if tmpFileErr != nil {
		panic(tmpFileErr)
	}

	var download = &Download{
		URL:      ts.URL,
		Filename: file.Name(),
	}

	if err := download.Load(); err != http.ErrMissingFile {
		t.Errorf("Load() should fail with %v, got %v instead", http.ErrMissingFile, err)
	}
}

func TestLoad(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.URL.Path)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	for _, c := range LoadCases {
		file, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
		defer os.Remove(file.Name())

		if tmpFileErr != nil {
			panic(tmpFileErr)
		}

		var download = &Download{
			URL:      ts.URL + c.word,
			Filename: file.Name(),
		}

		if err := download.Load(); err != nil {
			t.Errorf("Load() should not fail, got %v instead", err)
		}
	}
}
