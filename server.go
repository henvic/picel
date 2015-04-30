package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	VERSION                = "0.0.1"
	DEFAULT_ADDR           = ":8123"
	DEFAULT_BACKEND        = "http://localhost:8080/"
	SUCCESS_DECODE_MESSAGE = "Success. Image path parsed and decoded correctly"
	BAD_REQUEST_MESSAGE    = "Bad request."
)

var (
	addr        string
	backend     string
	verbose     bool
	flagVersion bool
	std         stdStruct
)

type stdStruct struct {
	out *log.Logger
	err *log.Logger
}

type Explain struct {
	Message    string    `json:"message"`
	Transform  Transform `json:"transform"`
	ErrorStack []string  `json:"errors"`
}

func init() {
	std.out = log.New(os.Stderr, "", log.LstdFlags)
	std.err = log.New(os.Stderr, "", log.LstdFlags)

	flag.StringVar(&addr, "addr", DEFAULT_ADDR, "Serving address")
	flag.StringVar(&backend, "storage", DEFAULT_BACKEND, "Image storage back-end server")
	flag.BoolVar(&verbose, "verbose", false, "Pipe image processing output to stderr/stdout")
	flag.BoolVar(&flagVersion, "version", false, "Print version information and quit")
}

func showVersion() {
	fmt.Println("ips version", VERSION)
}

func existsDependency(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func checkMissingDependencies(arg ...string) {
	var missing []string
	left := len(arg)

	for left > 0 {
		if !existsDependency(arg[left-1]) {
			missing = append(missing, arg[left-1])
		}
		left -= 1
	}

	if missing != nil {
		std.err.Println("Dependencies missing:", strings.Join(missing, ", "))
	}
}

func main() {
	flag.Parse()

	if flagVersion {
		showVersion()
		return
	}

	checkMissingDependencies("convert", "cwebp", "gif2webp")

	std.out.Println(fmt.Sprintf("Image Processing Service running on %v with backend %v", addr, backend))
	http.HandleFunc("/", handler)
	panic(http.ListenAndServe(addr, nil))
}

func buildExplain(transform Transform, err error, errs []error) Explain {
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
		Transform:  transform,
		ErrorStack: errorsMessages,
		Message:    message,
	}
}

func jsonEncodeTransformation(t Transform, errs []error, err error) string {
	res, _ := json.MarshalIndent(buildExplain(t, err, errs), "", "    ")

	return string(res)
}

func getOriginalUrl(image Image) string {
	name, _ := image.name()
	return backend + name
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

func processingHandler(filename string, t Transform, w http.ResponseWriter, r *http.Request) {
	if t.Raw {
		http.ServeFile(w, r, filename)
		return
	}

	output, _ := ioutil.TempFile(os.TempDir(), "ips")
	outputFilename := output.Name()
	defer os.Remove(outputFilename)

	err := Process(t, filename, output.Name())

	if err != nil {
		http.Error(w, "Processing error.", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, outputFilename)
}

func loadingHandler(t Transform, w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.TempFile(os.TempDir(), "ips")
	defer os.Remove(file.Name())
	filename := file.Name()

	url := getOriginalUrl(t.Image)

	_, err := Load(url, file.Name())

	if err != nil {
		http.NotFound(w, r)
		return
	}

	processingHandler(filename, t, w, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	transform, errs, err := Decode(r.URL.Path[1:], getDefaultRequestOutputFormat(r))

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
