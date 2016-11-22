/*
Package image provides encoding and processing for picel.
*/
package image

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// DefaultInputExtension is the default input format for images
	DefaultInputExtension = "jpg"

	// Raw is a special parameter to output a given image "as is" (acting like a proxy)
	Raw = "raw"
)

var (
	// ErrOffsetInvalid is returned when an offset is invalid
	ErrOffsetInvalid = errors.New("Offset is invalid")

	// ErrOffsetSeparator is returned when an off separator is not found
	ErrOffsetSeparator = errors.New("Offset separator not found")

	// ErrOffsetNonNegative is returned when an offset separator was given as a negative value
	ErrOffsetNonNegative = errors.New("x and y must be non-negative")

	// ErrBothDimensionEqualToZero is returned when the image size would be zero
	ErrBothDimensionEqualToZero = errors.New("At least x and y must be greater than zero")

	// ErrCropDimensionEqualToZero is returned when the image size would be zero after cropping
	ErrCropDimensionEqualToZero = errors.New("Both x and y must be greater than zero")

	// ErrNotCropFormat is returned when the crop format is invalid
	ErrNotCropFormat = errors.New("Not in crop format")

	// ErrInvalidCropDimensions is returned when the crop format dimensions is invalid
	ErrInvalidCropDimensions = errors.New("Invalid crop format dimensions")

	// ErrNonEmptyParameterQueue is returned when there are parameters left after processing a transformation
	ErrNonEmptyParameterQueue = errors.New("Can't process all parameters")
)

// Image structure
type Image struct {
	ID        string `json:"id"`
	Extension string `json:"extension"`
	Source    string `json:"source"`
}

// Crop parameters
type Crop struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Transform structure
type Transform struct {
	Image  `json:"image"`
	Path   string `json:"path"`
	Raw    bool   `json:"original"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Crop   Crop   `json:"crop"`
	Output string `json:"output"`
}

// Name of the image
func (i *Image) Name() (name string, fullname string) {
	fullname = i.ID

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

// Decode an image
func Decode(path string, defaultOutputFormat string) (transform Transform, errs []error, err error) {
	t := Transform{}

	t.Path = path

	paramsSubstringStart := getParamsSubstringStart(path)

	imgID := ""
	paramsString := ""
	output := ""

	if paramsSubstringStart == -1 {
		imgID, output = GetFilePathParts(path)
		t.Image.ID = UnescapePath(imgID)

		extension := output

		if extension == "" {
			extension = DefaultInputExtension
		}

		t.Image.Extension = UnescapePath(extension)
		t.Output = getOutputFormat(UnescapePath(output), defaultOutputFormat)
		_, fullname := t.Image.Name()
		t.Image.Source = fullname

		return t, errs, err
	}

	imgID = path[0 : paramsSubstringStart-1]
	paramsString, output = GetFilePathParts(path[paramsSubstringStart-1 : len(path)])
	t.Image.ID = UnescapePath(imgID)
	t.Output = getOutputFormat(UnescapePath(output), defaultOutputFormat)
	errs, err = extractParams(paramsString, UnescapePath(output), &t)
	_, fullname := t.Image.Name()
	t.Image.Source = fullname

	return t, errs, err
}

// Encode an image with a transformation
func Encode(transform Transform) (url string) {
	image := transform.Image
	url = EscapePath(image.ID)

	inputExtension := image.Extension

	if inputExtension == "" {
		inputExtension = DefaultInputExtension
	}

	if transform.Raw {
		url += "_" + Raw + "." + inputExtension

		return url
	}

	url += EncodeParam(encodeCrop(transform.Crop))

	url += EncodeParam(encodeDimension(transform.Width, transform.Height))

	if transform.Output != inputExtension && (inputExtension != DefaultInputExtension || transform.Output != "") {
		url += EncodeParam(EscapePath(inputExtension))
	}

	if transform.Output != "" {
		url += "." + EscapePath(transform.Output)
	}

	return url
}

// EscapePath of an image
func EscapePath(raw string) string {
	return strings.Replace(raw, "_", "__", -1)
}

// UnescapePath of an image
func UnescapePath(raw string) string {
	return strings.Replace(raw, "__", "_", -1)
}

func encodeCrop(c Crop) (crop string) {
	if c.Width != 0 && c.Height != 0 {
		crop = fmt.Sprintf("%dx%d:%dx%d", c.X, c.Y, c.Width, c.Height)
	}

	return crop
}

func encodeDimension(width int, height int) (dim string) {
	if width == 0 && height == 0 {
		return dim
	}

	if width > 0 {
		dim += fmt.Sprintf("%d", width)
	}

	dim += "x"

	if height > 0 {
		dim += fmt.Sprintf("%d", height)
	}

	return dim
}

// EncodeParam of an image
func EncodeParam(param string) string {
	if param != "" {
		param = "_" + param
	}

	return param
}

func getParamsSubstringStart(sp string) int {
	nextP := -1
	pivot := 0

	for {
		nextP = strings.Index(sp[pivot:], "_")

		if nextP == -1 {
			break
		}

		pivot += nextP + 1

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

func extractParams(part string, output string, t *Transform) (errs []error, err error) {
	params := strings.Split(part, "_")

	for i := range params {
		params[i] = UnescapePath(params[i])
	}

	pos := 1

	if len(params) == 2 && params[pos] == Raw {
		t.Raw = true
		t.Image.Extension = output
		return errs, err
	}

	crop, errsCrop := extractCrop(params[pos])

	if len(errsCrop) == 0 {
		t.Crop = crop
		pos++
	}

	errs = append(errs, errsCrop...)

	if pos < len(params) {
		width, height, errsResize := getDimensions(params[pos])

		if len(errsResize) == 0 {
			t.Width, t.Height = width, height
			pos++
		}

		errs = append(errs, errsResize...)
	}

	extension := output

	if pos != len(params) && params[pos] != "" {
		extension = params[pos]
		pos++
	}

	if extension == "" {
		extension = DefaultInputExtension
	}

	t.Image.Extension = UnescapePath(extension)

	if pos != len(params) {
		err = ErrNonEmptyParameterQueue
		errs = append(errs, err)
	}

	return errs, err
}

// GetFilePathParts separates extension from a given path
func GetFilePathParts(part string) (string, string) {
	last := strings.LastIndex(part, ".")

	if last == -1 {
		return part[0:len(part)], ""
	}

	return part[0:last], part[last+1 : len(part)]
}
