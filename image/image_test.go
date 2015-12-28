package image

import (
	"reflect"
	"testing"
)

type NameProvider struct {
	i        Image
	name     string
	fullname string
}

type EscapePathProvider struct {
	unescaped string
	escaped   string
}

type EncodeCropProvider struct {
	in   Crop
	want string
}

type EncodeDimensionProvider struct {
	in   Transform
	want string
}

type EncodeParamProvider struct {
	in   string
	want string
}

type ExtractCropProvider struct {
	in   string
	want Crop
}

type ExtractCropFailureProvider struct {
	in string
}

type GetParamsSubstringProvider struct {
	in   string
	want int
}

type GetOffsetsProvider struct {
	in string
	x  int
	y  int
}

type GetOffsetsFailureProvider struct {
	in string
}

type GetDimensionsProvider struct {
	in string
	x  int
	y  int
}

type GetDimensionsFailureProvider struct {
	in string
}

type GetOutputProvider struct {
	in     string
	prefix string
	suffix string
}

type DecodingFailureUnknownParameterProvider struct {
	in string
}

type DecodingFailureProvider struct {
	in string
}

type CompleteEncodingAndDecodingProvider struct {
	object Transform
	url    string
}

type DecodingToDefaultOutputFormatProvider struct {
	object Transform
	url    string
}

type IncompleteEncodingProvider struct {
	object Transform
	url    string
}

func TestName(t *testing.T) {
	for _, c := range NameCases {
		image := c.i
		name, fullname := image.Name()

		if name != c.name || fullname != c.fullname {
			t.Errorf("i.name() == %q %q, want %q %q", name, fullname, c.name, c.fullname)
		}
	}
}

func TestEscapePath(t *testing.T) {
	for _, c := range EscapePathCases {
		esc := EscapePath(c.unescaped)
		unesc := UnescapePath(c.escaped)

		if esc != c.escaped {
			t.Errorf("EscapePath(%q) == %q, want %q", c.unescaped, esc, c.escaped)
		}

		if unesc != c.unescaped {
			t.Errorf("EscapePath(%q) == %q, want %q", c.escaped, unesc, c.unescaped)
		}
	}
}

func TestEncodeCrop(t *testing.T) {
	for _, c := range EncodeCropCases {
		got := encodeCrop(c.in)

		if got != c.want {
			t.Errorf("encodeCrop(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEncodeDimension(t *testing.T) {
	for _, c := range EncodeDimensionCases {
		in := c.in
		got := encodeDimension(in.Width, in.Height)

		if got != c.want {
			t.Errorf("encodeDimension(%v, %v) == %v, want %v", in.Width, in.Height, got, c.want)
		}
	}
}

func TestEncodeParam(t *testing.T) {
	for _, c := range EncodeParamCases {
		got := EncodeParam(c.in)

		if got != c.want {
			t.Errorf("EncodeParam(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestExtractCrop(t *testing.T) {
	for _, c := range ExtractCropCases {
		crop, _ := extractCrop(c.in)

		if reflect.DeepEqual(crop, c.want) != true {
			t.Errorf("extractCrop(%v) == %v, want %v", c.in, crop, c.want)
		}
	}
}

func TestExtractCropFailure(t *testing.T) {
	for _, c := range ExtractCropFailureCases {
		_, err := extractCrop(c.in)

		if err == nil {
			t.Errorf("extractCrop(%q) should fail", c.in)
		}
	}
}

func TestGetParamsSubstringStart(t *testing.T) {
	for _, c := range GetParamsSubstringStartCases {
		got := getParamsSubstringStart(c.in)

		if got != c.want {
			t.Errorf("getParamsSubstringStart(%q) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestGetOffsets(t *testing.T) {
	for _, c := range GetOffsetsCases {
		x, y, err := getOffsets(c.in)

		if x != c.x || y != c.y || len(err) != 0 {
			t.Errorf("getOffsets(%q) == %dx%d, want %dx%d", c.in, x, y, c.x, c.y)
		}
	}
}

func TestGetOffsetsFailure(t *testing.T) {
	for _, c := range GetOffsetsFailureCases {
		_, _, err := getOffsets(c.in)

		if err == nil {
			t.Errorf("getOffsets(%q) should fail", c.in)
		}
	}
}

func TestGetDimensions(t *testing.T) {
	for _, c := range GetDimensionsCases {
		x, y, err := getDimensions(c.in)

		if x != c.x || y != c.y || len(err) != 0 {
			t.Errorf("getDimensions(%q) == %dx%d, want %dx%d", c.in, x, y, c.x, c.y)
		}
	}
}

func TestGetDimensionsFailure(t *testing.T) {
	for _, c := range GetDimensionsFailureCases {
		_, _, err := getDimensions(c.in)

		if err == nil {
			t.Errorf("getDimensions(%q) should fail", c.in)
		}
	}
}

func TestGetOutput(t *testing.T) {
	for _, c := range GetOutputCases {
		prefix, suffix := GetFilePathParts(c.in)

		if prefix != c.prefix || suffix != c.suffix {
			t.Errorf("GetFilePathParts(%q) == %q %q, want %q %q", c.in, prefix, suffix, c.prefix, c.suffix)
		}
	}
}

func TestDecodingFailureUnknownParameter(t *testing.T) {
	for _, c := range DecodingFailureUnknownParameterCases {
		_, _, err := Decode(c.in, "webp")

		if err == nil {
			t.Errorf("There should be errors due to invalid params for Decode(%q)", c.in)
		}
	}
}

func TestDecodingFailure(t *testing.T) {
	for _, c := range DecodingFailureCases {
		_, _, err := Decode(c.in, "jpg")

		if err != ErrNonEmptyParameterQueue {
			t.Errorf("There should be errors due to invalid / unprocessed params for Decode(%q)", c.in)
		}
	}
}

func TestCompleteEncodingAndDecoding(t *testing.T) {
	for _, c := range CompleteEncodingAndDecodingCases {
		gotUrl := Encode(c.object)

		if gotUrl != c.url {
			t.Errorf("Encode(%+v) == %v, want %v", c.object, gotUrl, c.url)
		}

		gotObject, _, err := Decode(c.url, "")

		if err != nil {
			t.Errorf("There should be no errors for Decode(%v)", c.url)
		}

		if reflect.DeepEqual(gotObject, c.object) != true {
			t.Errorf("Decode(%v) == %+v, want %+v", c.url, gotObject, c.object)
		}
	}
}

func TestDecodingToDefaultOutputFormat(t *testing.T) {
	for _, c := range DecodingToDefaultOutputFormatCases {
		gotObject, _, err := Decode(c.url, "other")

		if err != nil {
			t.Errorf("There should be no errors for Decode(%v)", c.url)
			t.Errorf("Decode(%v) == %+v, want %+v", c.url, gotObject, c.object)
		}

		if reflect.DeepEqual(gotObject, c.object) != true {
			t.Errorf("Decode(%v) == %+v, want %+v", c.url, gotObject, c.object)
		}
	}
}

func TestIncompleteEncoding(t *testing.T) {
	for _, c := range IncompleteEncodingCases {
		gotUrl := Encode(c.object)

		if gotUrl != c.url {
			t.Errorf("Encode(%+v) == %v, want %v", c.object, gotUrl, c.url)
		}
	}
}
