package image

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/henvic/picel/logger"
)

type ProcessProvider struct {
	input string
	t     Transform
}

type InvalidProcessProvider struct {
	t      Transform
	input  string
	output string
	err    error
}

func init() {
	// binary test assets are stored in a helper branch for neatness
	branch := exec.Command("git", "branch", "test_assets", "--track", "origin/test_assets", "-f")
	branch.Stderr = os.Stderr
	branch.Run()

	checkout := exec.Command("git", "checkout", "test_assets", "--", "../test_assets")
	checkout.Stderr = os.Stderr
	checkout.Run()

	gitRmCached := exec.Command("git", "rm", "--cached", "-r", "../test_assets")
	gitRmCached.Stderr = os.Stderr
	gitRmCached.Run()
}

func TestProcessInputFileNotFound(t *testing.T) {
	output, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")

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
		t.Errorf("Process(%q, %v) should fail", file, transform)
	}
}

func TestInvalidProcess(t *testing.T) {
	for _, c := range InvalidProcessCases {
		err := Process(c.t, "../"+c.input, c.output)

		if err != c.err {
			t.Errorf("Process(%v, %q, %q) should fail with %q, not %q", c.t, c.input, c.output, c.err, err)
		}
	}
}

func TestProcess(t *testing.T) {
	for _, c := range ProcessCases {
		output, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
		defer os.Remove(output.Name())

		if tmpFileErr != nil {
			panic(tmpFileErr)
		}

		err := Process(c.t, "../"+c.input, output.Name())

		if err != nil {
			t.Errorf("Process(%q, %v, %q) should not fail", "../"+c.input, c.t, output.Name())
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
	// don't run in parallel due to mocking logger.Stdout / logger.Stderr
	for _, c := range ProcessCasesForVerboseOn {
		output, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
		defer os.Remove(output.Name())

		if tmpFileErr != nil {
			panic(tmpFileErr)
		}

		var StdoutMock bytes.Buffer
		var StderrMock bytes.Buffer

		defaultStdout := logger.Stdout
		defaultStderr := logger.Stderr
		logger.Stdout = log.New(&StdoutMock, "", -1)
		logger.Stderr = log.New(&StderrMock, "", -1)
		Verbose = true
		err := Process(c.t, "../"+c.input, output.Name())
		Verbose = false
		logger.Stdout = defaultStdout
		logger.Stderr = defaultStderr

		if err != nil {
			t.Errorf("Process(%q, %v, %q) should not fail", "../"+c.input, c.t, output.Name())
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

		// convert uses Stderr in a strange way
		// http://www.imagemagick.org/discourse-server/viewtopic.php?t=9292
		if len(outMessages)+len(errMessages) == 100 {
			t.Errorf("Stderr / Stdout unexpectedly low")
		}
	}
}

func TestProcessFailureForEmptyFileWithVerboseOn(t *testing.T) {
	// don't run in parallel due to mocking Stdout / Stderr
	for _, c := range ProcessFailureForEmptyFileWithVerboseOnCases {
		output, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
		defer os.Remove(output.Name())

		if tmpFileErr != nil {
			panic(tmpFileErr)
		}

		var StdoutMock bytes.Buffer
		var StderrMock bytes.Buffer

		defaultStdout := logger.Stdout
		defaultStderr := logger.Stderr
		logger.Stdout = log.New(&StdoutMock, "", -1)
		logger.Stderr = log.New(&StderrMock, "", -1)
		Verbose = true
		err := Process(c.t, "../"+c.input, output.Name())
		Verbose = false
		logger.Stdout = defaultStdout
		logger.Stderr = defaultStderr

		if err == nil {
			t.Errorf("Process(%q, %v, %q) should fail", "../"+c.input, c.t, output.Name())
		}

		outMessages := StdoutMock.String()
		errMessages := StderrMock.String()

		// convert uses Stderr in a strange way
		// http://www.imagemagick.org/discourse-server/viewtopic.php?t=9292
		if len(outMessages)+len(errMessages) == 100 {
			t.Errorf("Stderr / Stdout unexpectedly low")
		}
	}
}
