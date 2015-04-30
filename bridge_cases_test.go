package main

var ProcessCases = []ProcessProvider{
	{"test_assets/golden-gate-bridge.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/golden-gate-bridge",
				Extension: "jpg",
			},
			Output: "jpg",
		}},
	{"test_assets/insects-2.JPEG",
		Transform{
			Image: Image{
				Id:        "test_assets/insects",
				Extension: "JPEG",
			},
			Output: "JPEG",
		}},
	{"test_assets/raccoons.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/raccoons",
				Extension: "jpg",
			},
			Width:  100,
			Output: "jpg",
		}},
	{"test_assets/golden-gate-bridge.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/golden-gate-bridge",
				Extension: "jpg",
			},
			Height: 100,
			Output: "jpg",
		}},
	{"test_assets/golden-gate-bridge.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/golden-gate-bridge",
				Extension: "jpg",
			},
			Width:  100,
			Height: 100,
			Output: "jpg",
		}},
	{"test_assets/golden-gate-bridge.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/golden-gate-bridge",
				Extension: "jpg",
			},
			Crop: Crop{
				X:      0,
				Y:      0,
				Width:  100,
				Height: 200,
			},
			Output: "jpg",
		}},
	{"test_assets/rocks_waves_big_sur_2.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/rocks_waves_big_sur_2",
				Extension: "jpg",
			},
			Output: "webp",
		}},
	{"test_assets/rocks_waves_big_sur_1.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/rocks_waves_big_sur_1",
				Extension: "jpg",
			},
			Width:  100,
			Output: "webp",
		}},
	{"test_assets/rocks_waves_big_sur_1.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/rocks_waves_big_sur_1",
				Extension: "jpg",
			},
			Height: 100,
			Output: "webp",
		}},
	{"test_assets/rocks_waves_big_sur_1.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/rocks_waves_big_sur_1",
				Extension: "jpg",
			},
			Width:  100,
			Height: 100,
			Output: "webp",
		}},
	{"test_assets/raccoons.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/raccoons",
				Extension: "jpg",
			},
			Crop: Crop{
				X:      0,
				Y:      0,
				Width:  100,
				Height: 200,
			},
			Output: "webp",
		}},
	{"test_assets/barter.gif",
		Transform{
			Image: Image{
				Id:        "test_assets/barter",
				Extension: "gif",
			},
			Crop: Crop{
				X:      0,
				Y:      0,
				Width:  100,
				Height: 200,
			},
			Output: "webp",
		}},
}

var ProcessCasesForVerboseOn = []ProcessProvider{
	{"test_assets/insects.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/insects",
				Extension: "jpg",
			},
			Output: "jpg",
		}},
	{"test_assets/insects.jpg",
		Transform{
			Image: Image{
				Id:        "test_assets/insects",
				Extension: "jpg",
			},
			Output: "webp",
		}},
	{"test_assets/barter.gif",
		Transform{
			Image: Image{
				Id:        "test_assets/barter",
				Extension: "gif",
			},
			Output: "webp",
		}},
	{"test_assets/barter.gif",
		Transform{
			Image: Image{
				Id:        "test_assets/barter",
				Extension: "gif",
			},
			Width:  100,
			Output: "gif",
		}},
	{"test_assets/barter.gif",
		Transform{
			Image: Image{
				Id:        "test_assets/barter",
				Extension: "gif",
			},
			Width: 100,
			Crop: Crop{
				X:      0,
				Y:      0,
				Width:  200,
				Height: 300,
			},
			Output: "webp",
		}},
}

var ProcessFailureForEmptyFileWithVerboseOnCases = []ProcessProvider{
	{"test_assets/empty_file.jpg",
		Transform{
			Image: Image{
				Id:        "empty_file",
				Extension: "jpg",
			},
			Output: "jpg",
		}},
	{"test_assets/empty_file.gif",
		Transform{
			Image: Image{
				Id:        "empty_file",
				Extension: "gif",
			},
			Height: 100,
			Output: "webp",
		}},
	{"test_assets/empty_file.jpg",
		Transform{
			Image: Image{
				Id:        "empty_file",
				Extension: "jpg",
			},
			Output: "webp",
		}},
}

var InvalidProcessCases = []InvalidProcessProvider{
	{Transform{
		Image: Image{
			Id:        "20120528-IMG_5236",
			Extension: "jpg",
		},
		Output: "unknown",
	},
		"foo.png",
		"foo.unknown"},
}
