package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultInputExtension = "jpg"
	defaultOutput         = ""
	RAW                   = "raw"
)

type Image struct {
	Id        string `json:"id"`
	Extension string `json:"extension"`
}

type Crop struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Transform struct {
	Image  `json:"image"`
	Raw    bool   `json:"original"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Crop   Crop   `json:"crop"`
	Output string `json:"output"`
}

func Decode(path string) (transform Transform, err error) {
	t := Transform{}

	paramsSubstringStart := getParamsSubstringStart(path)

	imgId := ""
	paramsString := ""
	output := ""

	if paramsSubstringStart == -1 {
		imgId, output = getFilePathParts(path)
		t.Image.Id = unescapeRawUrlParts(imgId)
		t.Image.Extension = unescapeRawUrlParts(output)
		t.Output = unescapeRawUrlParts(output)

		return t, err
	}

	imgId = path[0 : paramsSubstringStart-1]
	paramsString, output = getFilePathParts(path[paramsSubstringStart-1 : len(path)])
	err = extractParams(paramsString, output, &t)
	t.Image.Id = unescapeRawUrlParts(imgId)
	t.Output = unescapeRawUrlParts(output)

	return t, err
}

func Encode(transform Transform) (url string) {
	image := transform.Image
	url = escapeRawUrlParts(image.Id)

	inputExtension := image.Extension

	if inputExtension == "" {
		inputExtension = defaultInputExtension
	}

	if transform.Raw {
		url += "_" + RAW + "." + inputExtension

		return url
	}

	url += encodeParam(encodeCrop(transform.Crop))

	url += encodeParam(EncodeDimension(transform))

	if transform.Output != inputExtension {
		url += encodeParam(escapeRawUrlParts(inputExtension))
	}

	if transform.Output != "" {
		url += "." + escapeRawUrlParts(transform.Output)
	}

	return url
}

func escapeRawUrlParts(raw string) string {
	return strings.Replace(raw, "_", "__", -1)
}

func unescapeRawUrlParts(raw string) string {
	return strings.Replace(raw, "__", "_", -1)
}

func encodeCrop(c Crop) (crop string) {
	if c.Width != 0 && c.Height != 0 {
		crop = fmt.Sprintf("%dx%d:%dx%d", c.X, c.Y, c.Width, c.Height)
	}

	return crop
}

func EncodeDimension(transform Transform) (dim string) {
	if transform.Width == 0 && transform.Height == 0 {
		return dim
	}

	if transform.Width > 0 {
		dim += fmt.Sprintf("%d", transform.Width)
	}

	dim += "x"

	if transform.Height > 0 {
		dim += fmt.Sprintf("%d", transform.Height)
	}

	return dim
}

func encodeParam(param string) string {
	if param != "" {
		param = "_" + param
	}

	return param
}

func getParamsSubstringStart(sp string) int {
	next_ := -1
	pivot := 0

	for {
		next_ = strings.Index(sp[pivot:], "_")

		if next_ == -1 {
			break
		}

		pivot += next_ + 1

		if sp[pivot:][:1] != "_" && (pivot < 2 || sp[pivot-2:][:1] != "_") {
			return pivot
		}
	}

	return -1
}

func getDimensions(c string) (x int, y int, err error) {
	div := strings.Index(c, "x")

	if div == -1 {
		return x, y, errors.New("Dimensions separator not found")
	}

	current := c[0:div]

	if current != "" {
		x, err = strconv.Atoi(current)
	}

	if err != nil {
		return x, y, err
	}

	current = c[div+1 : len(c)]

	if current != "" {
		y, err = strconv.Atoi(current)
	}

	if err != nil {
		return x, y, err
	}

	if x < 0 || y < 0 {
		err = errors.New("x and y must be non-negative")
	}

	if x == 0 && y == 0 {
		err = errors.New("At least x and y must be greater than zero")
	}

	return x, y, err
}

func getCropDimensions(c string) (x int, y int, err error) {
	x, y, err = getDimensions(c)

	if x == 0 || y == 0 {
		err = errors.New("Both x and y must be greater than zero")
	}

	return x, y, err
}

func extractCrop(c string) (crop Crop, err error) {
	dot := strings.Index(c, ":")

	if dot == -1 {
		err = errors.New("Not in crop format")
		return crop, err
	}

	x, y, err1 := getDimensions(c[0:dot])
	width, height, err2 := getCropDimensions(c[dot+1 : len(c)])

	if err1 != nil || err2 != nil {
		err = errors.New("Invalid crop format dimensions")
	}

	crop = Crop{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}

	return crop, err
}

func extractParams(part string, output string, t *Transform) (err error) {
	params := strings.Split(part, "_")

	for i := range params {
		params[i] = unescapeRawUrlParts(params[i])
	}

	pos := 1

	if len(params) == 2 && params[pos] == RAW {
		t.Raw = true
		t.Image.Extension = output
		return err
	}

	crop, errCrop := extractCrop(params[pos])

	if errCrop == nil {
		t.Crop = crop
		pos += 1
	}

	width, height, errResize := getDimensions(params[pos])

	if errResize == nil {
		t.Width, t.Height = width, height
		pos += 1
	}

	extension := t.Output

	if params[pos] != t.Output {
		extension = params[pos]
		pos += 1
	}

	t.Image.Extension = extension

	if pos != len(params) {
		err = errors.New("Can't process all parameters")
	}

	return err
}

func getFilePathParts(part string) (string, string) {
	last := strings.LastIndex(part, ".")

	if last == -1 {
		return part[0:len(part)], defaultOutput
	}

	return part[0:last], part[last+1 : len(part)]
}
