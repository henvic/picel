/*
Package server provides server encoding, decoding and request handler for picel.
*/
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/henvic/picel/client"
	"github.com/henvic/picel/image"
)

const (
	SuccessDecodeMessage = "Success. Image path parsed and decoded correctly"
	BadRequestMessage    = "Bad request."
	HTTPSchema           = "http://"
	HTTPSSchema          = "https://"
	FlagHTTPSSchema      = "s:"
)

var (
	Backend         string
	Verbose         bool
	ErrMissingPath  = "Missing path"
	DownloadTimeout time.Duration
)

type Explain struct {
	Message    string          `json:"message"`
	Path       string          `json:"path"`
	Transform  image.Transform `json:"transform"`
	ErrorStack []string        `json:"errors"`
}

type crop struct {
	X      json.Number `json:"x"`
	Y      json.Number `json:"y"`
	Width  json.Number `json:"width"`
	Height json.Number `json:"height"`
}

type publicImage struct {
	Backend string      `json:"backend"`
	Path    string      `json:"path"`
	Raw     bool        `json:"raw"`
	Crop    crop        `json:"crop"`
	Width   json.Number `json:"width"`
	Height  json.Number `json:"height"`
	Output  string      `json:"output"`
}

func compressHost(raw string) string {
	if strings.Index(raw, HTTPSSchema) == 0 {
		return strings.Replace(raw, HTTPSSchema, FlagHTTPSSchema, 1)
	}

	return strings.Replace(raw, HTTPSchema, "", 1)
}

func expandHost(raw string) (source string) {
	https := false
	source = raw

	if strings.Index(source, FlagHTTPSSchema) == 0 {
		https = true
		source = strings.TrimPrefix(source, FlagHTTPSSchema)
	}

	switch https {
	case true:
		source = HTTPSSchema + source
	default:
		source = HTTPSchema + source
	}

	return source
}

func Decode(rawurl string, defaultOutputFormat string) (transform image.Transform, errs []error, err error) {
	rawurlIndex := strings.Index(rawurl, "/")

	host := rawurl
	path := ""

	if rawurlIndex != -1 {
		host = rawurl[0:rawurlIndex]
		path = rawurl[rawurlIndex+1:]
	}

	host = expandHost(host)

	transform, errs, err = image.Decode(path, defaultOutputFormat)

	_, fullname := transform.Image.Name()
	transform.Image.Source = host + "/" + fullname

	return transform, errs, err
}

func Encode(transform image.Transform) (url string) {
	url = image.Encode(transform)

	if Backend != "" {
		return compressHost(Backend) + "/" + url
	}

	source := transform.Image.Source
	_, fullname := transform.Image.Name()

	return compressHost(source[0:len(source)-len(fullname)]) + url
}

func buildExplain(path string, transform image.Transform, err error, errs []error) Explain {
	var errorsMessages []string
	var message string

	for i := range errs {
		errorsMessages = append(errorsMessages, fmt.Sprintf("%v", errs[i]))
	}

	if err != nil {
		message = fmt.Sprintf("%v", err)
	}

	if message == "" {
		message = SuccessDecodeMessage
	}

	return Explain{
		Message:    message,
		Path:       path,
		Transform:  transform,
		ErrorStack: errorsMessages,
	}
}

func jsonEncodeTransformation(path string, t image.Transform, errs []error, err error) string {
	res, _ := json.MarshalIndent(buildExplain(path, t, err, errs), "", "    ")

	return string(res)
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

	err := image.Process(t, filename, outputFilename)

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

	download := &client.Download{
		URL:      t.Image.Source,
		Filename: file.Name(),
	}

	if DownloadTimeout > 0*time.Second {
		download.Timeout(DownloadTimeout)
	}

	err := download.Load()

	if err != nil {
		http.NotFound(w, r)
		return
	}

	processingHandler(filename, t, w, r)
}

func encodeCrop(c crop) (param string) {
	if len(c.Width) != 0 && len(c.Height) != 0 {
		param = fmt.Sprintf("%sx%s:%sx%s", c.X, c.Y, c.Width, c.Height)
	}

	return param
}

func encodeDimension(width string, height string) (dim string) {
	if len(width) == 0 && len(height) == 0 {
		return dim
	}

	if len(width) != 0 {
		dim += fmt.Sprintf("%s", width)
	}

	dim += "x"

	if len(height) != 0 {
		dim += fmt.Sprintf("%s", height)
	}

	return dim
}

func createRequestPath(body io.Reader) (path string, err error) {
	decoder := json.NewDecoder(body)

	var pi publicImage
	var params []string

	err = decoder.Decode(&pi)

	if len(pi.Backend) != 0 {
		path = "/" + strings.TrimSuffix(compressHost(pi.Backend), "/")
	}

	id, extension := image.GetFilePathParts(pi.Path)

	id = strings.TrimPrefix(image.EscapePath(id), "/")

	if len(id) != 0 {
		path += "/" + id
	}

	if pi.Raw {
		path += "_" + image.RAW + "." + extension
		return path, err
	}

	params = append(params, encodeCrop(pi.Crop))
	params = append(params, encodeDimension(string(pi.Width), string(pi.Height)))

	if pi.Output != extension && (extension != image.DefaultInputExtension || len(pi.Output) != 0) {
		params = append(params, image.EscapePath(extension))
	}

	for _, param := range params {
		path += image.EncodeParam(param)
	}

	if len(pi.Output) != 0 {
		path += "." + image.EscapePath(pi.Output)
	}

	if pi.Path == "" {
		err = errors.New(ErrMissingPath)
	}

	return path, err
}

func prepare(r *http.Request) (transform image.Transform, reqPath string, errs []error, err error) {
	path := r.URL.Path[1:]
	reqPath = path
	var errRequestPath error

	if len(path) == 0 {
		var p string
		p, errRequestPath = createRequestPath(r.Body)

		if errRequestPath == nil && len(p) != 0 {
			path = p[1:]
			reqPath = path
		}
	}

	if Backend != "" {
		path = compressHost(Backend) + "/" + path
	}

	transform, errsDecode, err := Decode(path, getDefaultRequestOutputFormat(r))

	if errRequestPath != nil {
		errs = append(errs, errRequestPath)
		err = errRequestPath
	}

	errs = append(errs, errsDecode...)

	return transform, reqPath, errs, err
}

func Handler(w http.ResponseWriter, r *http.Request) {
	transform, path, errs, err := prepare(r)

	if r.URL.Query()["explain"] != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, jsonEncodeTransformation("/"+path, transform, errs, err))
		return
	}

	if err != nil {
		http.Error(w, BadRequestMessage, http.StatusBadRequest)
		return
	}

	loadingHandler(transform, w, r)
}
