/*
picel (picture element) is an image processing micro service.

	https://github.com/henvic/picel

*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/henvic/picel/image"
	"github.com/henvic/picel/logger"
	"github.com/henvic/picel/server"
)

const (
	// Version of latest picel release
	Version        = "0.0.1"
	defaultAddr    = ":8123"
	defaultBackend = ""
)

var (
	addr        string
	verbose     bool
	flagVersion bool
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "Serving address")
	flag.StringVar(&server.Backend, "backend", defaultBackend, "Image storage back-end server")
	flag.BoolVar(&verbose, "verbose", false, "Pipe image processing output to stderr/stdout")
	flag.BoolVar(&flagVersion, "version", false, "Print version information and quit")
}

func showVersion() {
	fmt.Println("picel version", Version)
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
		left--
	}

	if missing != nil {
		logger.Stderr.Println("Dependencies missing:", strings.Join(missing, ", "))
	}
}

func main() {
	flag.Parse()

	image.Verbose = verbose
	server.Verbose = verbose

	if flagVersion {
		showVersion()
		return
	}

	checkMissingDependencies("convert", "cwebp", "gif2webp")

	logger.Stdout.Println(fmt.Sprintf("picel started listening on %v", addr))

	if server.Backend != "" {
		logger.Stdout.Println(fmt.Sprintf("Single backend mode: %v", server.Backend))
	}

	http.HandleFunc("/", server.Handler)
	panic(http.ListenAndServe(addr, nil))
}
