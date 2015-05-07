package server

import (
	"bytes"
	"fmt"
	"github.com/henvic/picel/image"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type HostProvider struct {
	compressed string
	expanded   string
}

type GoodRequestProvider struct {
	url            string
	outputFiletype string
	meta           []string
	width          int
	height         int
	acceptWebp     bool
}

type BuildExplainProvider struct {
	transform image.Transform
	err       error
	errs      []error
	explain   Explain
}

type ServerProcessingFailureProvider struct {
	url string
}

type EncodingAndDecodingProvider struct {
	object image.Transform
	url    string
}

type EncodingAndDecodingForExplicitBackendProvider struct {
	object  image.Transform
	url     string
	backend string
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

func TestCompressAndExpandHost(t *testing.T) {
	t.Parallel()
	for _, c := range HostCases {
		compressed := compressHost(c.expanded)
		if compressed != c.compressed {
			t.Errorf("compressHost(%v) == %v, want %v", c.expanded, compressed, c.compressed)
		}

		expanded := expandHost(c.compressed)
		if expanded != c.expanded {
			t.Errorf("compressHost(%v) == %v, want %v", c.expanded, compressed, c.compressed)
		}
	}
}

func TestCompleteEncodingAndDecoding(t *testing.T) {
	t.Parallel()
	for _, c := range EncodingAndDecodingCases {
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

func TestEncodingForExplicitBackend(t *testing.T) {
	defaultBackend := Backend

	for _, c := range EncodingAndDecodingForExplicitBackendCases {
		Backend = c.backend
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

	Backend = defaultBackend
}

func TestBuildExplain(t *testing.T) {
	t.Parallel()
	for _, c := range BuildExplainCases {
		got := buildExplain(c.transform, c.err, c.errs)

		if reflect.DeepEqual(got, c.explain) != true {
			t.Errorf("buildExplain(%v, %v, %v) == %+v, want %+v", c.transform, c.err, c.errs, got, c.explain)
		}
	}
}

func TestJsonEncodeTransformation(t *testing.T) {
	t.Parallel()
	path := "s:example.net/foo_137x0:737x450_800x600_jpg.webp"
	reference := "../explain_example.json"

	content, err := ioutil.ReadFile(reference)

	if err != nil {
		panic(err)
	}

	want := string(content)

	actual := jsonEncodeTransformation(Decode(path, "foo"))
	jsonEncodeTransformation(image.Transform{}, nil, nil)

	if actual != want {
		t.Errorf("Expected JSON for %v doesn't match with the result saved as %v", path, reference)
	}
}

func TestServerExplain(t *testing.T) {
	t.Parallel()
	url := "/s:example.net/foo_137x0:737x450_800x600_jpg.webp"
	reference := "../explain_example.json"

	req, _ := http.NewRequest("GET", url+"?explain", nil)
	w := httptest.NewRecorder()
	http.HandlerFunc(Handler).ServeHTTP(w, req)

	content, err := ioutil.ReadFile(reference)

	if err != nil {
		panic(err)
	}

	want := string(content)

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}

	if w.Body.String() != want {
		t.Errorf("Expected JSON for %v doesn't match with the result saved as %v", url, reference)
	}
}

func TestServerSingleBackendExplain(t *testing.T) {
	// don't run in parallel due to mocking Backend
	url := "/foo_137x0:737x450_800x600_jpg.webp"
	reference := "../explain_example.json"

	defaultBackend := Backend
	Backend = "https://example.net"
	req, _ := http.NewRequest("GET", url+"?explain", nil)
	w := httptest.NewRecorder()
	http.HandlerFunc(Handler).ServeHTTP(w, req)
	Backend = defaultBackend

	content, err := ioutil.ReadFile(reference)

	if err != nil {
		panic(err)
	}

	want := string(content)

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}

	if w.Body.String() != want {
		t.Errorf("Expected JSON for %v doesn't match with the result saved as %v", url, reference)
	}
}

func TestServerBadRequest(t *testing.T) {
	// don't run in parallel due to mocking Backend
	url := "/_"

	defaultBackend := Backend
	Backend = "https://localhost/"
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	http.HandlerFunc(Handler).ServeHTTP(w, req)
	Backend = defaultBackend

	if w.Code != http.StatusBadRequest {
		t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusBadRequest)
	}

	if w.Body.String() != BAD_REQUEST_MESSAGE+"\n" {
		t.Errorf("Bad request body message response is %v, want %v", w.Body.String(), BAD_REQUEST_MESSAGE)
	}
}

func TestServerNotFound(t *testing.T) {
	t.Parallel()
	url := "/not-found_640x"

	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	http.HandlerFunc(Handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusNotFound)
	}
}

func identifyImageDetails(filename string, meta []string, transform image.Transform, t *testing.T) {
	// imagick doesn't support decoding some standards like .gif
	if meta == nil {
		return
	}

	identify := exec.Command("identify", "-verbose", filename)
	out, identifyImageErr := identify.CombinedOutput()
	info := string(out)

	if identifyImageErr != nil {
		t.Errorf("Error while identifying image: %v", identifyImageErr)
	}

	for _, v := range meta {
		if strings.LastIndex(info, v) == -1 {
			t.Errorf("Error identifying image: want %v, but it was not found for %+v", v, transform)
		}
	}
}

func compareImage(output string, transform image.Transform, width int, height int, t *testing.T) {
	reference := "../test_assets/compare" + transform.Path

	compare := exec.Command("compare", "-metric", "AE", "-fuzz", "15%", reference, output, "/dev/null")

	// always ignore compare exit code
	out, _ := compare.CombinedOutput()
	content := string(out)

	if strings.Index(content, "is not supported.") != -1 || strings.Index(content, "no decode") != -1 {
		fmt.Println(fmt.Sprintf("Jumping comparing for %v: no support / decoder for ImageMagick compare tool", reference))
		return
	}

	lastLine := strings.LastIndex(content, "\n")

	if lastLine != -1 {
		content = content[lastLine+1 : len(content)]
	}

	abs, err := strconv.Atoi(content)

	if err != nil {
		panic(err)
	}

	diff := float64(abs) / float64(width*height)

	if diff > 0.15 {
		t.Errorf("Image %+v is very different than expected", transform)
	}
}

func verifyGoodRequest(compBackend string, c GoodRequestProvider, t *testing.T) {
	url := "/" + compBackend + c.url
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	defaultOutput := "jpg"

	if c.acceptWebp {
		req.Header.Set("Accept", "image/webp,*/*;q=0.8")
		defaultOutput = "webp"
	}

	http.HandlerFunc(Handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusOK)
	}

	actualContentType := w.Header().Get("Content-Type")

	if actualContentType != c.outputFiletype {
		t.Errorf("Content-Type is %v, want %v", actualContentType, c.outputFiletype)
	}

	// this checking is very unsafe as can lead to false positives
	// but works with the current test cases
	transform, _, _ := image.Decode(c.url, defaultOutput)

	if transform.Raw {
		_, filename := transform.Image.Name()
		reference, _ := ioutil.ReadFile("../test_assets/" + filename)

		if bytes.Compare(reference, w.Body.Bytes()) != 0 {
			t.Errorf("Raw file for %v differ from what is expected", filename)
		}

		return
	}

	file, tmpFileErr := ioutil.TempFile(os.TempDir(), "picel")
	filename := file.Name()
	defer os.Remove(filename)
	defer file.Close()

	if tmpFileErr != nil {
		panic(tmpFileErr)
	}

	io.Copy(file, w.Body)

	identifyImageDetails(filename, c.meta, transform, t)
	compareImage(filename, transform, c.width, c.height, t)
}

func TestServerGoodRequests(t *testing.T) {
	t.Parallel()
	fsHandler := func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "../", "", -1)
		http.ServeFile(w, r, "../test_assets/"+path[1:])
	}

	ts := httptest.NewServer(http.HandlerFunc(fsHandler))
	defer ts.Close()

	compBackend := compressHost(ts.URL)

	for _, c := range GoodRequestsCases {
		verifyGoodRequest(compBackend, c, t)
	}
}

func TestServerProcessingFailure(t *testing.T) {
	t.Parallel()
	fsHandler := func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "../", "", -1)
		http.ServeFile(w, r, "../test_assets/"+path[1:])
	}

	ts := httptest.NewServer(http.HandlerFunc(fsHandler))
	defer ts.Close()

	Backend = ts.URL + "/"

	for _, c := range ServerProcessingFailureCases {
		req, _ := http.NewRequest("GET", c.url, nil)
		w := httptest.NewRecorder()

		http.HandlerFunc(Handler).ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusInternalServerError)
		}
	}
}
