package main

import (
	"reflect"
	"testing"
)

func TestName(t *testing.T) {
	cases := []struct {
		i        Image
		name     string
		fullname string
	}{
		{Image{
			Id:        "help/staff",
			Extension: "jpg",
		}, "staff.jpg", "help/staff.jpg"},
		{Image{
			Id:        "section/help/staff",
			Extension: "jpg",
		}, "staff.jpg", "section/help/staff.jpg"},
		{Image{
			Id:        "dog",
			Extension: "png",
		}, "dog.png", "dog.png"},
		{Image{
			Id:        "dog",
			Extension: "",
		}, "dog", "dog"},
	}
	for _, c := range cases {
		name, fullname := c.i.name()

		if name != c.name || fullname != c.fullname {
			t.Errorf("c.i.name() == %q %q, want %q %q", name, fullname, c.name, c.fullname)
		}
	}
}

func TestEscapeRawUrlParts(t *testing.T) {
	cases := []struct {
		unescaped string
		escaped   string
	}{
		{"", ""},
		{"_", "__"},
		{"__", "____"},
		{"x_", "x__"},
		{"_y", "__y"},
		{"x_y", "x__y"},
	}
	for _, c := range cases {
		esc := escapeRawUrlParts(c.unescaped)
		unesc := unescapeRawUrlParts(c.escaped)

		if esc != c.escaped {
			t.Errorf("escapeRawUrlParts(%q) == %q, want %q", c.unescaped, esc, c.escaped)
		}

		if unesc != c.unescaped {
			t.Errorf("escapeRawUrlParts(%q) == %q, want %q", c.escaped, unesc, c.unescaped)
		}
	}
}

func TestEncodeCrop(t *testing.T) {
	cases := []struct {
		in   Crop
		want string
	}{
		{Crop{
			X:      0,
			Y:      0,
			Width:  10,
			Height: 10,
		}, "0x0:10x10"},
	}
	for _, c := range cases {
		got := encodeCrop(c.in)

		if got != c.want {
			t.Errorf("encodeCrop(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEncodeDimension(t *testing.T) {
	cases := []struct {
		in   Transform
		want string
	}{
		{Transform{
			Width:  0,
			Height: 0,
		}, ""},
		{Transform{
			Width: 10,
		}, "10x"},
		{Transform{
			Height: 10,
		}, "x10"},
		{Transform{
			Width:  10,
			Height: 10,
		}, "10x10"},
	}
	for _, c := range cases {
		got := EncodeDimension(c.in)

		if got != c.want {
			t.Errorf("EncodeDimension(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestEncodeParam(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"x", "_x"},
	}
	for _, c := range cases {
		got := encodeParam(c.in)

		if got != c.want {
			t.Errorf("encodeParam(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestExtractCrop(t *testing.T) {
	cases := []struct {
		in   string
		want Crop
	}{
		{"10x20:400x300",
			Crop{
				X:      10,
				Y:      20,
				Width:  400,
				Height: 300,
			}},
	}
	for _, c := range cases {
		crop, _ := extractCrop(c.in)

		if reflect.DeepEqual(crop, c.want) != true {
			t.Errorf("extractCrop(%v) == %v, want %v", c.in, crop, c.want)
		}
	}
}

func TestExtractCropFailure(t *testing.T) {
	cases := []struct {
		in string
	}{
		{""},
		{"10x20:x300"},
		{"10x20:400x"},
		{"10x:x300"},
		{"20:400"},
	}
	for _, c := range cases {
		_, err := extractCrop(c.in)

		if err == nil {
			t.Errorf("extractCrop(%q) should fail", c.in)
		}
	}
}

func TestGetParamsSubstringStart(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", -1},
		{"little__kittens", -1},
		{"dogs_4x4.png", 5},
		{"animals/turtles__newborn_4x4.jpg", 25},
	}
	for _, c := range cases {
		got := getParamsSubstringStart(c.in)

		if got != c.want {
			t.Errorf("getParamsSubstringStart(%q) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestGetDimensions(t *testing.T) {
	cases := []struct {
		in string
		x  int
		y  int
	}{
		{"500x100", 500, 100},
		{"300x", 300, 0},
		{"x300", 0, 300},
	}
	for _, c := range cases {
		x, y, err := getDimensions(c.in)

		if x != c.x || y != c.y || len(err) != 0 {
			t.Errorf("getDimensions(%q) == %dx%d, want %dx%d", c.in, x, y, c.x, c.y)
		}
	}
}

func TestGetDimensionsFailure(t *testing.T) {
	cases := []struct {
		in string
	}{
		{""},
		{"-1x10"},
		{"10x-1"},
		{"-1x-1"},
		{"x"},
		{"yx10"},
		{"10xy"},
		{"yxy"},
	}
	for _, c := range cases {
		_, _, err := getDimensions(c.in)

		if err == nil {
			t.Errorf("getDimensions(%q) should fail", c.in)
		}
	}
}

func TestGetOutput(t *testing.T) {
	cases := []struct {
		in     string
		prefix string
		suffix string
	}{
		{"", "", ""},
		{"file", "file", ""},
		{"file.out", "file", "out"},
	}
	for _, c := range cases {
		prefix, suffix := getFilePathParts(c.in)

		if prefix != c.prefix || suffix != c.suffix {
			t.Errorf("getFilePathParts(%q) == %q %q, want %q %q", c.in, prefix, suffix, c.prefix, c.suffix)
		}
	}
}

func TestDecodingFailureUnknownParameter(t *testing.T) {
	cases := []struct {
		in string
	}{
		{"la__office/newborn__bunnies_raw_stars.jpg"},
	}
	for _, c := range cases {
		_, err, _ := Decode(c.in)

		if err == nil {
			t.Errorf("There should be errors due to invalid params for Decode(%q)", c.in)
		}
	}
}

func TestDecodingFailure(t *testing.T) {
	cases := []struct {
		in string
	}{
		{"la__office/newborn__bunnies_.jpg"},
		{"la__office/newborn__bunnies_400x200:300_gif.jpg"},
		{"la__office/newborn__bunnies_400x200:nox300_gif.jpg"},
		{"la__office/newborn__bunnies_400x200:300xno_gif.jpg"},
	}
	for _, c := range cases {
		_, err, _ := Decode(c.in)

		if err == nil {
			t.Errorf("There should be errors due to invalid params for Decode(%q)", c.in)
		}
	}
}

func TestCompleteEncodingAndDecoding(t *testing.T) {
	cases := []struct {
		object Transform
		url    string
	}{
		{Transform{
			Image: Image{
				Id:        "help/staff",
				Extension: "jpg",
			},
			Output: "jpg",
		}, "help/staff.jpg"},
		{Transform{
			Image: Image{
				Id:        "help/staff",
				Extension: "jpg",
			},
			Output: "webp",
		}, "help/staff_jpg.webp"},
		{Transform{
			Image: Image{
				Id:        "airplane_flying_low",
				Extension: "jpg",
			},
			Output: "webp",
		}, "airplane__flying__low_jpg.webp"},
		{Transform{
			Image: Image{
				Id:        "help/foo",
				Extension: "jpg",
			},
			Output: "",
			Width:  400,
			Height: 800,
		}, "help/foo_400x800_jpg"},
		{Transform{
			Image: Image{
				Id:        "help/foo",
				Extension: "jpg",
			},
			Output: "",
			Width:  400,
		}, "help/foo_400x_jpg"},
		{Transform{
			Image: Image{
				Id:        "help/foo",
				Extension: "jpg",
			},
			Output: "",
			Height: 800,
		}, "help/foo_x800_jpg"},
		{Transform{
			Image: Image{
				Id:        "adoption_shelters_in_nyc/pretty_dogs",
				Extension: "jpg",
			},
			Output: "webp",
			Width:  400,
			Height: 800,
		}, "adoption__shelters__in__nyc/pretty__dogs_400x800_jpg.webp"},
		{Transform{
			Image: Image{
				Id:        "airplane_360",
				Extension: "gif",
			},
			Output: "gif",
		}, "airplane__360.gif"},
		{Transform{
			Image: Image{
				Id:        "airplane_360",
				Extension: "gif",
			},
		}, "airplane__360_gif"},
		{Transform{
			Image: Image{
				Id:        "airplane_360",
				Extension: "gif",
			},
			Output: "webp",
		}, "airplane__360_gif.webp"},
		{Transform{
			Image: Image{
				Id:        "foo",
				Extension: "jpg",
			},
			Output: "",
		}, "foo_jpg"},
		{Transform{
			Image: Image{
				Id:        "foo",
				Extension: "jpg",
			},
			Width:  800,
			Height: 600,
			Crop: Crop{
				X:      137,
				Y:      0,
				Width:  737,
				Height: 450,
			},
			Output: "webp",
		}, "foo_137x0:737x450_800x600_jpg.webp",
		},
		{Transform{
			Image: Image{
				Id:        "adoption_shelters_in_nyc/pretty_dogs",
				Extension: "jpg",
			},
			Width:  800,
			Height: 600,
			Crop: Crop{
				X:      137,
				Y:      1,
				Width:  737,
				Height: 451,
			},
			Output: "webp",
		}, "adoption__shelters__in__nyc/pretty__dogs_137x1:737x451_800x600_jpg.webp",
		},
		{Transform{
			Image: Image{
				Id:        "adoption_shelters_in_nyc/pretty_dogs",
				Extension: "jpg",
			},
			Crop: Crop{
				X:      137,
				Y:      1,
				Width:  737,
				Height: 451,
			},
			Output: "webp",
		}, "adoption__shelters__in__nyc/pretty__dogs_137x1:737x451_jpg.webp",
		},
		{Transform{
			Image: Image{
				Id:        "la_office/newborn_bunnies",
				Extension: "jpg",
			},
			Raw:    true,
			Output: "jpg",
		}, "la__office/newborn__bunnies_raw.jpg",
		},
	}
	for _, c := range cases {
		gotUrl := Encode(c.object)

		if gotUrl != c.url {
			t.Errorf("Encode(%+v) == %v, want %v", c.object, gotUrl, c.url)
		}

		gotObject, err, _ := Decode(c.url)

		if err != nil {
			t.Errorf("There should be no errors for Decode(%v)", c.url)
		}

		if reflect.DeepEqual(gotObject, c.object) != true {
			t.Errorf("Decode(%v) == %+v, want %+v", c.url, gotObject, c.object)
		}
	}
}

func TestIncompleteEncoding(t *testing.T) {
	cases := []struct {
		object Transform
		url    string
	}{
		{Transform{
			Image: Image{
				Id:        "la_office/newborn_bunnies",
				Extension: "jpg",
			},
			Raw: true,
		}, "la__office/newborn__bunnies_raw.jpg",
		},
		{Transform{
			Image: Image{
				Id: "foo",
			},
			Output: "",
		}, "foo_jpg"},
		{Transform{
			Image: Image{
				Id: "help/staff",
			},
			Output: "webp",
		}, "help/staff_jpg.webp"},
	}
	for _, c := range cases {
		gotUrl := Encode(c.object)

		if gotUrl != c.url {
			t.Errorf("Encode(%+v) == %v, want %v", c.object, gotUrl, c.url)
		}
	}
}
