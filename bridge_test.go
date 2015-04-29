package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
)

func init() {
	// binary test assets are stored in a helper branch for neatness
	exec.Command("git", "checkout", "test_assets", "--", "test_assets").Run()
	exec.Command("git", "rm", "--cached", "-r", "test_assets").Run()
}

func TestProcessInputFileNotFound(t *testing.T) {
	t.Parallel()
	output, tmpFileErr := ioutil.TempFile(os.TempDir(), "ips")

	if tmpFileErr != nil {
		panic(tmpFileErr)
	}

	defer os.Remove(output.Name())

	file := "not-found"

	transform := Transform{
		Image: Image{
			Id:        "20120528-IMG_5236",
			Extension: "jpg",
		},
		Output: "jpg",
	}

	err := Process(transform, file, output.Name())

	if err == nil {
		t.Errorf("Process(%q, %q) should fail", file, transform)
	}
}

func TestInvalidProcess(t *testing.T) {
	t.Parallel()
	cases := []struct {
		t      Transform
		input  string
		output string
	}{
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
	for _, c := range cases {
		err := Process(c.t, c.input, c.output)

		if err != ErrOutputFormatNotSupported {
			t.Errorf("Process(%q, %q, %q) unknown output format should make it fail", c.t, c.input, c.output)
		}
	}
}

func TestProcess(t *testing.T) {
	t.Parallel()
	cases := []struct {
		filename string
		t        Transform
	}{
		{"test_assets/golden-gate-bridge.jpg",
			Transform{
				Image: Image{
					Id:        "test_assets/golden-gate-bridge",
					Extension: "jpg",
				},
				Output: "jpg",
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
	for _, c := range cases {
		output, tmpFileErr := ioutil.TempFile(os.TempDir(), "ips")
		defer os.Remove(output.Name())

		if tmpFileErr != nil {
			panic(tmpFileErr)
		}

		err := Process(c.t, c.filename, output.Name())

		if err != nil {
			t.Errorf("Process(%q, %q, %q) should not fail", c.filename, c.t, output.Name())
		}

		fileInfo, fileInfoErr := os.Stat(output.Name())

		if fileInfoErr != nil {
			panic(fileInfoErr)
		}

		if fileInfo.Size() == 0 {
			t.Errorf("Processed file size is zero")
		}
	}
}

func TestProcessWithVerboseOn(t *testing.T) {
	// don't run in parallel due to mocking std.out / std.err
	cases := []struct {
		filename string
		t        Transform
	}{
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
	for _, c := range cases {
		output, tmpFileErr := ioutil.TempFile(os.TempDir(), "ips")
		defer os.Remove(output.Name())

		if tmpFileErr != nil {
			panic(tmpFileErr)
		}

		var StdoutMock bytes.Buffer
		var StderrMock bytes.Buffer

		defaultStdout, defaultStderr := std.out, std.err
		std.out = log.New(&StdoutMock, "", -1)
		std.err = log.New(&StderrMock, "", -1)
		verbose = true
		err := Process(c.t, c.filename, output.Name())
		verbose = false
		std.out = defaultStdout
		std.err = defaultStderr

		if err != nil {
			t.Errorf("Process(%q, %q, %q) should not fail", c.filename, c.t, output.Name())
		}

		fileInfo, fileInfoErr := os.Stat(output.Name())

		if fileInfoErr != nil {
			panic(fileInfoErr)
		}

		if fileInfo.Size() == 0 {
			t.Errorf("Processed file size is zero")
		}

		outMessages := StdoutMock.String()
		errMessages := StderrMock.String()

		// convert uses stderr in a strange way
		// http://www.imagemagick.org/discourse-server/viewtopic.php?t=9292
		if len(outMessages)+len(errMessages) == 100 {
			t.Errorf("Stderr / stdout unexpectedly low")
		}
	}
}

func TestProcessFailureForEmptyFileWithVerboseOn(t *testing.T) {
	// don't run in parallel due to mocking std.out / std.err
	cases := []struct {
		filename string
		t        Transform
	}{
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
	for _, c := range cases {
		output, tmpFileErr := ioutil.TempFile(os.TempDir(), "ips")
		defer os.Remove(output.Name())

		if tmpFileErr != nil {
			panic(tmpFileErr)
		}

		var StdoutMock bytes.Buffer
		var StderrMock bytes.Buffer

		defaultStdout, defaultStderr := std.out, std.err
		std.out = log.New(&StdoutMock, "", -1)
		std.err = log.New(&StderrMock, "", -1)
		verbose = true
		err := Process(c.t, c.filename, output.Name())
		verbose = false
		std.out = defaultStdout
		std.err = defaultStderr

		if err == nil {
			t.Errorf("Process(%q, %q, %q) should fail", c.filename, c.t, output.Name())
		}

		outMessages := StdoutMock.String()
		errMessages := StderrMock.String()

		// convert uses stderr in a strange way
		// http://www.imagemagick.org/discourse-server/viewtopic.php?t=9292
		if len(outMessages)+len(errMessages) == 100 {
			t.Errorf("Stderr / stdout unexpectedly low")
		}
	}
}
