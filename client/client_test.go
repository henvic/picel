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
	_, err := Load("0/foo.png", "")

	if err == nil {
		t.Errorf("Load(0/foo.png, ) should fail")
	}
}

func TestLoadFromInvalidSchema(t *testing.T) {
	file, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
	defer os.Remove(file.Name())

	if tmpFileErr != nil {
		panic(tmpFileErr)
	}

	_, err := Load("0/foo.png", file.Name())

	if err == nil {
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

	size, err := Load(ts.URL, file.Name())

	if err != ErrBackend {
		t.Errorf("Load(%v, %v) == %v %v should fail", ts.URL, file.Name(), size, err)
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

	size, err := Load(ts.URL, file.Name())

	if err != http.ErrMissingFile {
		t.Errorf("Load(%v, %v) == %v %v should fail", ts.URL, file.Name(), size, err)
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

		size, err := Load(ts.URL+c.word, file.Name())

		fileInfo, fileInfoErr := os.Stat(file.Name())

		if fileInfoErr != nil {
			panic(fileInfoErr)
		}

		if err != nil || int(size) != len(c.word) || fileInfo.Size() != size {
			t.Errorf("Load(%q, %q) == %q %q, want %q %v", ts.URL, file.Name(), size, err, len(c.word), nil)
		}
	}
}
