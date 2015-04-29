package main

var NameCases = []NameProvider{
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

var EscapeRawUrlPartsCases = []EscapeRawUrlPartsProvider{
	{"", ""},
	{"_", "__"},
	{"__", "____"},
	{"x_", "x__"},
	{"_y", "__y"},
	{"x_y", "x__y"},
}

var EncodeCropCases = []EncodeCropProvider{
	{Crop{
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
	}, "0x0:10x10"},
}

var EncodeDimensionCases = []EncodeDimensionProvider{
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

var EncodeParamCases = []EncodeParamProvider{
	{"", ""},
	{"x", "_x"},
}

var ExtractCropCases = []ExtractCropProvider{
	{"10x20:400x300",
		Crop{
			X:      10,
			Y:      20,
			Width:  400,
			Height: 300,
		}},
}

var ExtractCropFailureCases = []ExtractCropFailureProvider{
	{""},
	{"10x20:x300"},
	{"10x20:400x"},
	{"10x:x300"},
	{"20:400"},
}

var GetParamsSubstringStartCases = []GetParamsSubstringProvider{
	{"", -1},
	{"little__kittens", -1},
	{"dogs_4x4.png", 5},
	{"animals/turtles__newborn_4x4.jpg", 25},
}

var GetOffsetsCases = []GetOffsetsProvider{
	{"500x100", 500, 100},
	{"300x", 300, 0},
	{"x300", 0, 300},
	{"0x0", 0, 0},
}

var GetOffsetsFailureCases = []GetOffsetsFailureProvider{
	{""},
	{"-1x10"},
	{"10x-1"},
	{"-1x-1"},
	{"x"},
	{"yx10"},
	{"10xy"},
	{"yxy"},
}

var GetDimensionsCases = []GetDimensionsProvider{
	{"500x100", 500, 100},
	{"300x", 300, 0},
	{"x300", 0, 300},
}

var GetDimensionsFailureCases = []GetDimensionsFailureProvider{
	{""},
	{"-1x10"},
	{"10x-1"},
	{"-1x-1"},
	{"x"},
	{"yx10"},
	{"10xy"},
	{"yxy"},
	{"0x0"},
}

var GetOutputCases = []GetOutputProvider{
	{"", "", ""},
	{"file", "file", ""},
	{"file.out", "file", "out"},
}

var DecodingFailureUnknownParameterCases = []DecodingFailureUnknownParameterProvider{
	{"la__office/newborn__bunnies_raw_stars.jpg"},
}

var DecodingFailureCases = []DecodingFailureProvider{
	{"_"},
	{"la__office/newborn__bunnies_.jpg"},
	{"la__office/newborn__bunnies_400x200:300_gif.jpg"},
	{"la__office/newborn__bunnies_400x200:nox300_gif.jpg"},
	{"la__office/newborn__bunnies_400x200:300xno_gif.jpg"},
}

var CompeteEncodingAndDecodingCases = []CompeteEncodingAndDecodingProvider{
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
			Extension: "webp",
		},
		Output: "webp",
	}, "help/staff.webp"},
	{Transform{
		Image: Image{
			Id:        "help/staff",
			Extension: "webp",
		},
		Width:  800,
		Output: "webp",
	}, "help/staff_800x.webp"},
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
			Id:        "dog",
			Extension: "jpg",
		},
		Output: "",
	}, "dog"},
	{Transform{
		Image: Image{
			Id:        "help/foo",
			Extension: "jpg",
		},
		Output: "",
		Width:  400,
		Height: 800,
	}, "help/foo_400x800"},
	{Transform{
		Image: Image{
			Id:        "help/foo",
			Extension: "jpg",
		},
		Output: "",
		Width:  400,
	}, "help/foo_400x"},
	{Transform{
		Image: Image{
			Id:        "help/foo",
			Extension: "jpg",
		},
		Output: "",
		Height: 800,
	}, "help/foo_x800"},
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
	}, "foo"},
	{Transform{
		Image: Image{
			Id:        "foo_bah_h",
			Extension: "jpg",
		},
		Crop: Crop{
			X:      0,
			Y:      0,
			Width:  800,
			Height: 400,
		},
		Output: "jpg",
	}, "foo__bah__h_0x0:800x400.jpg"},
	{Transform{
		Image: Image{
			Id:        "foo_bah_h",
			Extension: "jpg",
		},
		Crop: Crop{
			X:      300,
			Y:      300,
			Width:  800,
			Height: 400,
		},
		Output: "jpg",
	}, "foo__bah__h_300x300:800x400.jpg"},
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

var DecodingToDefaultOutputFormatCases = []DecodingToDefaultOutputFormatProvider{
	{Transform{
		Image: Image{
			Id:        "la_office/newborn_bunnies",
			Extension: "jpg",
		},
		Raw:    false,
		Output: "other",
	}, "la__office/newborn__bunnies",
	},
	{Transform{
		Image: Image{
			Id:        "la_office/newborn_bunnies",
			Extension: "jpg",
		},
		Raw:    false,
		Output: "other",
	}, "la__office/newborn__bunnies_jpg",
	},
	{Transform{
		Image: Image{
			Id:        "la_office/newborn_bunnies",
			Extension: "other",
		},
		Raw:    false,
		Output: "other",
	}, "la__office/newborn__bunnies_other",
	},
	{Transform{
		Image: Image{
			Id:        "la_office/newborn_bunnies",
			Extension: "gif",
		},
		Raw:    false,
		Output: "gif",
	}, "la__office/newborn__bunnies.gif",
	},
	{Transform{
		Image: Image{
			Id:        "big_sur",
			Extension: "jpg",
		},
		Crop: Crop{
			X:      0,
			Y:      0,
			Width:  600,
			Height: 600,
		},
		Output: "other",
	}, "big__sur_0x0:600x600",
	},
}

var IncompleteEncodingCases = []IncompleteEncodingProvider{
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
	}, "foo"},
	{Transform{
		Image: Image{
			Id: "help/staff",
		},
		Output: "webp",
	}, "help/staff_jpg.webp"},
}
