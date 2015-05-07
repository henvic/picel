package server

import (
	"errors"
	"github.com/henvic/picel/image"
)

var HostCases = []HostProvider{
	{
		"google.com",
		"http://google.com",
	},
	{
		"s:google.com",
		"https://google.com",
	},
}

var EncodingAndDecodingCases = []EncodingAndDecodingProvider{
	{image.Transform{
		Image: image.Image{
			Id:        "help/staff",
			Extension: "jpg",
			Source:    "http://127.0.0.1/help/staff.jpg",
		},
		Path:   "help/staff.jpg",
		Output: "jpg",
	},
		"127.0.0.1/help/staff.jpg"},
	{image.Transform{
		Image: image.Image{
			Id:        "help/staff",
			Extension: "webp",
			Source:    "http://remote.local/help/staff.webp",
		},
		Path:   "help/staff.webp",
		Output: "webp",
	}, "remote.local/help/staff.webp"},
	{image.Transform{
		Image: image.Image{
			Id:        "help/staff",
			Extension: "webp",
			Source:    "https://localhost/help/staff.webp",
		},
		Path:   "help/staff_800x.webp",
		Width:  800,
		Output: "webp",
	}, "s:localhost/help/staff_800x.webp"},
}

var EncodingAndDecodingForExplicitBackendCases = []EncodingAndDecodingForExplicitBackendProvider{
	{image.Transform{
		Image: image.Image{
			Id:        "help/staff",
			Extension: "jpg",
			Source:    "http://127.0.0.1/help/staff.jpg",
		},
		Path:   "help/staff.jpg",
		Output: "jpg",
	},
		"127.0.0.1/help/staff.jpg",
		"127.0.0.1"},
	{image.Transform{
		Image: image.Image{
			Id:        "help/staff",
			Extension: "webp",
			Source:    "http://remote.local/help/staff.webp",
		},
		Path:   "help/staff.webp",
		Output: "webp",
	},
		"remote.local/help/staff.webp",
		"remote.local"},
	{image.Transform{
		Image: image.Image{
			Id:        "help/staff",
			Extension: "webp",
			Source:    "https://localhost/help/staff.webp",
		},
		Width:  800,
		Path:   "help/staff_800x.webp",
		Output: "webp",
	},
		"s:localhost/help/staff_800x.webp",
		"s:localhost"},
}

var GoodRequestsCases = []GoodRequestProvider{
	{
		"/rocks__waves__big__sur__2_raw.jpg",
		"image/jpeg",
		[]string{},
		true},
	{
		"/rocks__waves__big__sur__2_500x.jpg",
		"image/jpeg",
		[]string{
			"Format: JPEG",
			"Geometry: 500x333",
			"Interlace: None",
		},
		true},
	{
		"/insects-2_500x.JPEG",
		"image/jpeg",
		[]string{
			"Format: JPEG",
			"Geometry: 500x333",
			"Interlace: None",
		},
		true},
	{
		"/rocks__waves__big__sur__1_0x0:600x300_112x56",
		"image/jpeg",
		[]string{
			"Geometry: 112x56+0+0",
			"Interlace: None",
		},
		false},
	{
		"/big__sur_0x0:600x600_100x",
		"image/webp",
		[]string{
			"Geometry: 100x100+0+0",
			"Interlace: None",
		},
		true},
	{
		"/big__sur_0x0:600x600_100x",
		"image/jpeg",
		[]string{
			"Geometry: 100x100+0+0",
			"Interlace: None",
		},
		false},
	{
		"/additive-color_70x30:20x20.png",
		"image/png",
		[]string{
			"Geometry: 20x20+0+0",
			"red: 1-bit",
			"green: 1-bit",
			"blue: 1-bit",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Green:\n      min: 0 (0)\n      max: 0 (0)",
			"Blue:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_130x150:20x20.png",
		"image/png",
		[]string{
			"Geometry: 20x20+0+0",
			"red: 1-bit",
			"green: 1-bit",
			"blue: 1-bit",
			"Red:\n      min: 0 (0)\n      max: 0 (0)",
			"Green:\n      min: 0 (0)\n      max: 0 (0)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_50x150:20x20.png",
		"image/png",
		[]string{
			"Geometry: 20x20+0+0",
			"red: 1-bit",
			"green: 1-bit",
			"blue: 1-bit",
			"Red:\n      min: 0 (0)\n      max: 0 (0)",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_90x140:20x20.png",
		"image/png",
		[]string{
			"Geometry: 20x20+0+0",
			"red: 1-bit",
			"green: 1-bit",
			"blue: 1-bit",
			"Red:\n      min: 0 (0)\n      max: 0 (0)",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_120x70:20x20.png",
		"image/png",
		[]string{
			"Geometry: 20x20+0+0",
			"red: 1-bit",
			"green: 1-bit",
			"blue: 1-bit",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Green:\n      min: 0 (0)\n      max: 0 (0)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_40x70:20x20.png",
		"image/png",
		[]string{
			"Geometry: 20x20+0+0",
			"red: 1-bit",
			"green: 1-bit",
			"blue: 1-bit",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_0x0:20x20.png",
		"image/png",
		[]string{
			"Geometry: 20x20+0+0",
			"Colorspace: Gray",
			"Gray:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_70x30:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_130x150:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 0 (0)\n      max: 0 (0)",
			"Green:\n      min: 0 (0)\n      max: 0 (0)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_50x150:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 0 (0)\n      max: 0 (0)",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_90x140:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_120x70:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_40x70:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		true},
	{
		"/additive-color_0x0:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Colorspace: Gray",
			"Gray:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		true},
	{
		"/barter.gif",
		"image/gif",
		[]string{
			"Format: GIF",
			"Geometry: 400x300+0+0",
		},
		true},
	{
		"/barter_20x.gif",
		"image/gif",
		[]string{
			"Format: GIF",
			"Geometry: 20x15+0+0",
		},
		true},
	{
		"/barter_gif",
		"image/webp",
		nil,
		true},
	{
		"/barter_20x_gif",
		"image/webp",
		nil,
		true},
}

var BuildExplainCases = []BuildExplainProvider{
	{
		image.Transform{},
		nil,
		nil,
		Explain{
			Message:    SUCCESS_DECODE_MESSAGE,
			Transform:  image.Transform{},
			ErrorStack: nil},
	},
	{
		image.Transform{},
		errors.New("xyz"),
		[]error{errors.New("foo"), errors.New("testing")},
		Explain{
			Message:    "xyz",
			Transform:  image.Transform{},
			ErrorStack: []string{"foo", "testing"},
		}},
}

var ServerProcessingFailureCases = []ServerProcessingFailureProvider{
	{"/empty__file.jpg"},
	{"/insects_jpg.xoo"},
}
