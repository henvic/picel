package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	WEBP_QUALITY    = "92"
	IMAGICK_QUALITY = "92"
)

var (
	ErrOutputFormatNotSupported = errors.New("The requested output format is not supported")
)

func Process(t Transform, input string, output string) (err error) {
	formats := map[string]string{
		"jpg":  "Imagick",
		"gif":  "Imagick",
		"png":  "Imagick",
		"pdf":  "Imagick",
		"webp": "Cwebp",
	}

	tool, valid := formats[t.Output]

	if !valid {
		return ErrOutputFormatNotSupported
	}

	if tool == "Imagick" {
		return processImagick(t, input, output)
	}

	return processWebp(t, input, output)
}

func callProgram(name string, params []string) error {
	cmd := exec.Command(name, params...)
	var bOut bytes.Buffer
	var bErr bytes.Buffer
	cmd.Stdout = &bOut
	cmd.Stderr = &bErr

	if verbose {
		std.out.Println(fmt.Sprintf("%v %v", name, strings.Join(params, " ")))
	}

	cmdErr := cmd.Run()

	if verbose {
		std.out.Println(string(bOut.Bytes()))
		std.err.Println(string(bErr.Bytes()))
	}

	return cmdErr
}

func processWebp(t Transform, input string, output string) (err error) {
	if t.Extension != "gif" {
		return processCwebp(t, input, output)
	}

	if t.Crop.Width != 0 || t.Crop.Height != 0 || t.Width != 0 || t.Height != 0 {
		t.Output = "gif"
		err = processImagick(t, input, output)
		t.Output = "webp"
		input = output

		if err != nil {
			return err
		}
	}

	return processGif2Webp(input, output)
}

func processGif2Webp(input string, output string) (err error) {
	var params []string

	params = append(params, "-q")
	params = append(params, WEBP_QUALITY)

	if verbose {
		params = append(params, "-v")
	}

	params = append(params, input)
	params = append(params, "-o")
	params = append(params, output)

	return callProgram("gif2webp", params)
}

func processCwebp(t Transform, input string, output string) (err error) {
	var params []string

	params = append(params, "-q")
	params = append(params, WEBP_QUALITY)

	if t.Crop.Width != 0 && t.Crop.Height != 0 {
		params = append(params, "-crop")
		params = append(params, fmt.Sprintf("%d", t.Crop.X))
		params = append(params, fmt.Sprintf("%d", t.Crop.Y))
		params = append(params, fmt.Sprintf("%d", t.Crop.Width))
		params = append(params, fmt.Sprintf("%d", t.Crop.Height))
	}

	if t.Width != 0 || t.Height != 0 {
		params = append(params, "-resize")
		params = append(params, fmt.Sprintf("%d", t.Width))
		params = append(params, fmt.Sprintf("%d", t.Height))
	}

	if verbose {
		params = append(params, "-v")
	}

	params = append(params, input)
	params = append(params, "-o")
	params = append(params, output)

	return callProgram("cwebp", params)
}

func processImagick(t Transform, input string, output string) (err error) {
	var params []string

	if verbose {
		params = append(params, "-verbose")
	}

	params = append(params, "-quality")

	params = append(params, IMAGICK_QUALITY)

	params = append(params, input)

	params = append(params, "-strip")

	if t.Crop.Width != 0 && t.Crop.Height != 0 {
		params = append(params, "-crop")
		params = append(params, fmt.Sprintf("%dx%d+%d+%d", t.Crop.Width, t.Crop.Height, t.Crop.X, t.Crop.Y))
		params = append(params, "+repage")
	}

	if t.Width != 0 || t.Height != 0 {
		resize := ""

		if t.Width > 0 {
			resize += fmt.Sprintf("%d", t.Width)
		}

		resize += "x"

		if t.Height > 0 {
			resize += fmt.Sprintf("%d", t.Height)
		}

		params = append(params, "-resize")
		params = append(params, resize)
	}

	output = t.Output + ":" + output

	params = append(params, output)

	return callProgram("convert", params)
}
