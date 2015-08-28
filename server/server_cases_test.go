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
		`{
			"backend":
			"REPLACE_ON_TEST",
			"path": "rocks_waves_big_sur_2.jpg",
			"output": "jpg",
			"raw": true
		}`,
		"/rocks__waves__big__sur__2_raw.jpg",
		"image/jpeg",
		[]string{},
		1600,
		1067,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "rocks_waves_big_sur_2.jpg",
			"output": "jpg",
			"width": 500
		}`,
		"/rocks__waves__big__sur__2_500x.jpg",
		"image/jpeg",
		[]string{
			"Format: JPEG",
			"Geometry: 500x333",
			"Interlace: None",
		},
		500,
		333,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "insects-2.JPEG",
			"output": "JPEG",
			"width": "500"
		}`,
		"/insects-2_500x.JPEG",
		"image/jpeg",
		[]string{
			"Format: JPEG",
			"Geometry: 500x333",
			"Interlace: None",
		},
		500,
		333,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "rocks_waves_big_sur_1.jpg",
			"width": "112",
			"height": 56,
			"crop": {
				"x": 0,
				"y": "0",
				"width": 600,
				"height": 300
			}
		}`,
		"/rocks__waves__big__sur__1_0x0:600x300_112x56",
		"image/jpeg",
		[]string{
			"Geometry: 112x56+0+0",
			"Interlace: None",
		},
		112,
		56,
		false},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "big_sur.jpg",
			"width": 100,
			"crop": {
				"x": "0",
				"y": 0,
				"width": "600",
				"height": "600"
			}
		}`,
		"/big__sur_0x0:600x600_100x",
		"image/webp",
		[]string{
			"Geometry: 100x100+0+0",
			"Interlace: None",
		},
		100,
		100,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "big_sur.jpg",
			"width": 100,
			"crop": {
				"x": "0",
				"y": 0,
				"width": "600",
				"height": "600"
			}
		}`,
		"/big__sur_0x0:600x600_100x",
		"image/jpeg",
		[]string{
			"Geometry: 100x100+0+0",
			"Interlace: None",
		},
		100,
		100,
		false},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "/additive-color.png",
			"crop": {
				"x": "70",
				"y": "30",
				"width": "20",
				"height": "20"
			},
			"output": "png"
		}`,
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
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": "130",
				"y": "150",
				"width": "20",
				"height": "20"
			},
			"output": "png"
		}`,
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
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": "50",
				"y": "150",
				"width": "20",
				"height": "20"
			},
			"output": "png"
		}`,
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
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "/additive-color.png",
			"crop": {
				"x": "90",
				"y": "140",
				"width": "20",
				"height": "20"
			},
			"output": "png"
		}`,
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
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": "120",
				"y": "70",
				"width": "20",
				"height": "20"
			},
			"output": "png"
		}`,
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
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": 40,
				"y": 70,
				"width": 20,
				"height": 20
			},
			"output": "png"
		}`,
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
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": "0",
				"y": "0",
				"width": "20",
				"height": "20"
			},
			"output": "png"
		}`,
		"/additive-color_0x0:20x20.png",
		"image/png",
		[]string{
			"Geometry: 20x20+0+0",
			"Colorspace: Gray",
			"Gray:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": "70",
				"y": "30",
				"width": "20",
				"height": "20"
			},
			"output": "webp"
		}`,
		"/additive-color_70x30:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": "130",
				"y": "150",
				"width": "20",
				"height": "20"
			},
			"output": "webp"
		}`,
		"/additive-color_130x150:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 0 (0)\n      max: 0 (0)",
			"Green:\n      min: 0 (0)\n      max: 0 (0)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": "50",
				"y": "150",
				"width": "20",
				"height": "20"
			},
			"output": "webp"
		}`,
		"/additive-color_50x150:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 0 (0)\n      max: 0 (0)",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": 90,
				"y": 140,
				"width": 20,
				"height": 20
			},
			"output": "webp"
		}`,
		"/additive-color_90x140:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": "120",
				"y": "70",
				"width": "20",
				"height": "20"
			},
			"output": "webp"
		}`,
		"/additive-color_120x70:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 255 (1)\n      max: 255 (1)",
			"Interlace: None",
		},
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": 40,
				"y": 70,
				"width": 20,
				"height" : 20
			},
			"output": "webp"
		}`,
		"/additive-color_40x70:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Red:\n      min: 255 (1)\n      max: 255 (1)",
			"Green:\n      min: 255 (1)\n      max: 255 (1)",
			"Blue:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "additive-color.png",
			"crop": {
				"x": 0,
				"y": 0,
				"width": 20,
				"height" : 20
			},
			"output": "webp"
		}`,
		"/additive-color_0x0:20x20_png.webp",
		"image/webp",
		[]string{
			"Geometry: 20x20+0+0",
			"Colorspace: Gray",
			"Gray:\n      min: 0 (0)\n      max: 0 (0)",
			"Interlace: None",
		},
		20,
		20,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "barter.gif",
			"output": "gif"
		}`,
		"/barter.gif",
		"image/gif",
		[]string{
			"Format: GIF",
			"Geometry: 400x300+0+0",
		},
		400,
		300,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "barter.gif",
			"width": "20",
			"output": "gif"
		}`,
		"/barter_20x.gif",
		"image/gif",
		[]string{
			"Format: GIF",
			"Geometry: 20x15+0+0",
		},
		20,
		15,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "barter.gif"
		}`,
		"/barter_gif",
		"image/webp",
		nil,
		400,
		300,
		true},
	{
		`{
			"backend": "REPLACE_ON_TEST",
			"path": "barter.gif",
			"width": "20"
		}`,
		"/barter_20x_gif",
		"image/webp",
		nil,
		20,
		15,
		true},
}

var BadRequestsURLCases = []BadRequestProvider{
	{
		"/_",
	},
}

var BadRequestsJSONCases = []BadRequestProvider{
	{`{`}, {`{}`}, {`{"width": "200"}`},
}

var BuildExplainCases = []BuildExplainProvider{
	{
		"/xyz",
		image.Transform{},
		nil,
		nil,
		Explain{
			Message:    SuccessDecodeMessage,
			Path:       "/xyz",
			Transform:  image.Transform{},
			ErrorStack: nil},
	},
	{
		"/",
		image.Transform{},
		errors.New("xyz"),
		[]error{errors.New("foo"), errors.New("testing")},
		Explain{
			Message:    "xyz",
			Path:       "/",
			Transform:  image.Transform{},
			ErrorStack: []string{"foo", "testing"},
		}},
}

var ServerProcessingFailureCases = []ServerProcessingFailureProvider{
	{"/empty__file.jpg"},
	{"/insects_jpg.xoo"},
}

var EncodeCropCases = []EncodeCropProvider{
	{crop{
		X:      "1",
		Y:      "2",
		Width:  "3",
		Height: "4",
	}, "1x2:3x4"},
	{crop{}, ""},
}

var EncodeDimensionCases = []EncodeDimensionProvider{
	{publicImage{
		Width:  "",
		Height: "",
	}, ""},
	{publicImage{
		Width: "10",
	}, "10x"},
	{publicImage{
		Height: "10",
	}, "x10"},
	{publicImage{
		Width:  "10",
		Height: "10",
	}, "10x10"},
}

var CreateRequestPathCases = []CreateRequestPathProvider{
	{
		doc:  `{"path": "foo.jpg"}`,
		path: "/foo",
	}, {
		doc:  `{"path": "foo.jpg", "backend": "https://localhost/"}`,
		path: "/s:localhost/foo",
	}, {
		doc:  `{"path": "foo.gif", "raw": true}`,
		path: "/foo_raw.gif",
	}, {
		doc:  `{"path": "bah.jpg", "raw": true}`,
		path: "/bah_raw.jpg",
	}, {
		doc:  `{"path": "bah.jpg", "crop": {"x": 0, "y": 0, "width": 100, "height": 200}}`,
		path: "/bah_0x0:100x200",
	}, {
		doc:  `{"path": "bah.jpg", "crop": {"x": "0", "y": "0", "width": "100", "height": "200"}}`,
		path: "/bah_0x0:100x200",
	}, {
		doc:  `{"path": "bah.jpg", "width": 100}`,
		path: "/bah_100x",
	}, {
		doc:  `{"path": "bah.jpg", "width": "100"}`,
		path: "/bah_100x",
	}, {
		doc:  `{"path": "bah.jpg", "height": 100}`,
		path: "/bah_x100",
	}, {
		doc:  `{"path": "bah.jpg", "height": "100"}`,
		path: "/bah_x100",
	}, {
		doc:  `{"path": "bah.jpg", "width": 40, "height": "100"}`,
		path: "/bah_40x100",
	}, {
		doc:  `{"path": "bah.jpg", "width": "40", "height": 100}`,
		path: "/bah_40x100",
	}, {
		doc:  `{"path": "bah.gif", "width": "40", "output": "webp"}`,
		path: "/bah_40x_gif.webp",
	}, {
		doc:  `{"path": "foo_bah.jpg", "width": "40", "output": "jpg"}`,
		path: "/foo__bah_40x.jpg",
	},
}
