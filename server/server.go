/*
picel (picture element) is an image processing micro service.

	https://github.com/henvic/picel

*/

package server

import (
	"encoding/json"
	"fmt"
	"github.com/henvic/picel/client"
	"github.com/henvic/picel/image"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	SUCCESS_DECODE_MESSAGE = "Success. Image path parsed and decoded correctly"
	BAD_REQUEST_MESSAGE    = "Bad request."
)

var (
	Backend string
	Verbose bool
)

type Explain struct {
	Message    string          `json:"message"`
	Transform  image.Transform `json:"transform"`
	ErrorStack []string        `json:"errors"`
}

func buildExplain(transform image.Transform, err error, errs []error) Explain {
	var errorsMessages []string
	var message string

	for i := range errs {
		errorsMessages = append(errorsMessages, fmt.Sprintf("%v", errs[i]))
	}

	if err != nil {
		message = fmt.Sprintf("%v", err)
	}

	if message == "" {
		message = SUCCESS_DECODE_MESSAGE
	}

	return Explain{
		Message:    message,
		Transform:  transform,
		ErrorStack: errorsMessages,
	}
}

func jsonEncodeTransformation(t image.Transform, errs []error, err error) string {
	res, _ := json.MarshalIndent(buildExplain(t, err, errs), "", "    ")

	return string(res)
}

func getSourceUrl(image image.Image) string {
	name, _ := image.Name()
	return Backend + name
}

func isWebpCompatible(r *http.Request) bool {
	accept := r.Header["Accept"]
	return len(accept) != 0 && strings.Index(accept[0], "image/webp") != -1
}

func getDefaultRequestOutputFormat(r *http.Request) string {
	if isWebpCompatible(r) {
		return "webp"
	}

	return "jpg"
}

func processingHandler(filename string, t image.Transform, w http.ResponseWriter, r *http.Request) {
	if t.Raw {
		http.ServeFile(w, r, filename)
		return
	}

	output, _ := ioutil.TempFile(os.TempDir(), "picel")
	outputFilename := output.Name()
	defer os.Remove(outputFilename)

	err := image.Process(t, filename, output.Name())

	if err != nil {
		http.Error(w, "Processing error.", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, outputFilename)
}

func loadingHandler(t image.Transform, w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.TempFile(os.TempDir(), "picel")
	defer os.Remove(file.Name())
	filename := file.Name()

	url := getSourceUrl(t.Image)

	_, err := client.Load(url, file.Name())

	if err != nil {
		http.NotFound(w, r)
		return
	}

	processingHandler(filename, t, w, r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	transform, errs, err := image.Decode(r.URL.Path[1:], getDefaultRequestOutputFormat(r))

	if r.URL.Query()["explain"] != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, jsonEncodeTransformation(transform, errs, err))
		return
	}

	if err != nil {
		http.Error(w, BAD_REQUEST_MESSAGE, http.StatusBadRequest)
		return
	}

	loadingHandler(transform, w, r)
}
