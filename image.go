package picel

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrOffsetInvalid            = errors.New("Offset is invalid")
	ErrOffsetSeparator          = errors.New("Offset separator not found")
	ErrOffsetNonNegative        = errors.New("x and y must be non-negative")
	ErrBothDimensionEqualToZero = errors.New("At least x and y must be greater than zero")
	ErrCropDimensionEqualToZero = errors.New("Both x and y must be greater than zero")
	ErrNotCropFormat            = errors.New("Not in crop format")
	ErrInvalidCropDimensions    = errors.New("Invalid crop format dimensions")
	ErrNonEmptyParameterQueue   = errors.New("Can't process all parameters")
)

const (
	defaultInputExtension = "jpg"
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

func (i *Image) name() (name string, fullname string) {
	fullname = i.Id

	if i.Extension != "" {
		fullname += "." + i.Extension
	}

	last := strings.LastIndex(fullname, "/")

	if last == -1 {
		return fullname, fullname
	}

	return fullname[last+1 : len(fullname)], fullname
}

func getOutputFormat(output string, defaultOutputFormat string) string {
	if output == "" {
		return defaultOutputFormat
	}

	return output
}

func Decode(path string, defaultOutputFormat string) (transform Transform, errs []error, err error) {
	t := Transform{}

	paramsSubstringStart := getParamsSubstringStart(path)

	imgId := ""
	paramsString := ""
	output := ""

	if paramsSubstringStart == -1 {
		imgId, output = getFilePathParts(path)
		t.Image.Id = unescapeRawUrlParts(imgId)

		extension := output

		if extension == "" {
			extension = defaultInputExtension
		}

		t.Image.Extension = unescapeRawUrlParts(extension)
		t.Output = getOutputFormat(unescapeRawUrlParts(output), defaultOutputFormat)

		return t, errs, err
	}

	imgId = path[0 : paramsSubstringStart-1]
	paramsString, output = getFilePathParts(path[paramsSubstringStart-1 : len(path)])
	t.Image.Id = unescapeRawUrlParts(imgId)
	t.Output = getOutputFormat(unescapeRawUrlParts(output), defaultOutputFormat)
	err, errs = extractParams(paramsString, unescapeRawUrlParts(output), &t)

	return t, errs, err
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

	if transform.Output != inputExtension && (inputExtension != defaultInputExtension || transform.Output != "") {
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

		if len(sp[pivot:]) == 0 || (sp[pivot:][:1] != "_" && (pivot < 2 || sp[pivot-2:][:1] != "_")) {
			return pivot
		}
	}

	return -1
}

func getOffsets(c string) (x int, y int, errs []error) {
	var err error

	if len(c) <= 1 {
		errs = append(errs, ErrOffsetInvalid)
		return x, y, errs
	}

	div := strings.Index(c, "x")

	if div == -1 {
		errs = append(errs, ErrOffsetSeparator)
		return x, y, errs
	}

	current := c[0:div]

	if current != "" {
		x, err = strconv.Atoi(current)
	}

	if err != nil {
		errs = append(errs, err)
		return x, y, errs
	}

	current = c[div+1 : len(c)]

	if current != "" {
		y, err = strconv.Atoi(current)
	}

	if err != nil {
		errs = append(errs, err)
		return x, y, errs
	}

	if x < 0 || y < 0 {
		errs = append(errs, ErrOffsetNonNegative)
	}

	return x, y, errs
}

func getDimensions(c string) (x int, y int, errs []error) {
	x, y, errs = getOffsets(c)

	if x == 0 && y == 0 {
		errs = append(errs, ErrBothDimensionEqualToZero)
	}

	return x, y, errs
}

func getCropDimensions(c string) (x int, y int, errs []error) {
	x, y, errs = getDimensions(c)

	if x == 0 || y == 0 {
		errs = append(errs, ErrCropDimensionEqualToZero)
	}

	return x, y, errs
}

func extractCrop(c string) (crop Crop, errs []error) {
	dot := strings.Index(c, ":")

	if dot == -1 {
		errs = append(errs, ErrNotCropFormat)
		return crop, errs
	}

	x, y, errs1 := getOffsets(c[0:dot])
	width, height, errs2 := getCropDimensions(c[dot+1 : len(c)])

	errs = append(errs, errs1...)
	errs = append(errs, errs2...)

	if len(errs1) != 0 || len(errs2) != 0 {
		errs = append(errs, ErrInvalidCropDimensions)
	}

	crop = Crop{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}

	return crop, errs
}

func extractParams(part string, output string, t *Transform) (err error, errs []error) {
	params := strings.Split(part, "_")

	for i := range params {
		params[i] = unescapeRawUrlParts(params[i])
	}

	pos := 1

	if len(params) == 2 && params[pos] == RAW {
		t.Raw = true
		t.Image.Extension = output
		return err, errs
	}

	crop, errsCrop := extractCrop(params[pos])

	if len(errsCrop) == 0 {
		t.Crop = crop
		pos += 1
	}

	errs = append(errs, errsCrop...)

	if pos < len(params) {
		width, height, errsResize := getDimensions(params[pos])

		if len(errsResize) == 0 {
			t.Width, t.Height = width, height
			pos += 1
		}

		errs = append(errs, errsResize...)
	}

	extension := output

	if pos != len(params) && params[pos] != "" {
		extension = params[pos]
		pos += 1
	}

	if extension == "" {
		extension = defaultInputExtension
	}

	t.Image.Extension = unescapeRawUrlParts(extension)

	if pos != len(params) {
		err = ErrNonEmptyParameterQueue
		errs = append(errs, err)
	}

	return err, errs
}

func getFilePathParts(part string) (string, string) {
	last := strings.LastIndex(part, ".")

	if last == -1 {
		return part[0:len(part)], ""
	}

	return part[0:last], part[last+1 : len(part)]
}
