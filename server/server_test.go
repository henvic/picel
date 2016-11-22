package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/henvic/picel/image"
)

type HostProvider struct {
	compressed string
	expanded   string
}

type GoodRequestProvider struct {
	requestBody    string
	url            string
	outputFiletype string
	meta           []string
	width          int
	height         int
	acceptWebp     bool
}

type BadRequestProvider struct {
	request string
}

type BuildExplainProvider struct {
	path      string
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

type EncodeCropProvider struct {
	in   crop
	want string
}

type EncodeDimensionProvider struct {
	in   publicImage
	want string
}

type CreateRequestPathProvider struct {
	doc  string
	path string
}

var ts *httptest.Server

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

	fsHandler := func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "../", "", -1)
		http.ServeFile(w, r, "../test_assets/"+path[1:])
	}

	ts = httptest.NewServer(http.HandlerFunc(fsHandler))
}

func TestCompressAndExpandHost(t *testing.T) {
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
	for _, c := range EncodingAndDecodingCases {
		gotURL := Encode(c.object)

		if gotURL != c.url {
			t.Errorf("Encode(%+v) == %v, want %v", c.object, gotURL, c.url)
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
		gotURL := Encode(c.object)

		if gotURL != c.url {
			t.Errorf("Encode(%+v) == %v, want %v", c.object, gotURL, c.url)
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
	for _, c := range BuildExplainCases {
		got := buildExplain(c.path, c.transform, c.err, c.errs)

		if reflect.DeepEqual(got, c.explain) != true {
			t.Errorf("buildExplain(%v, %v, %v, %v) == %+v, want %+v", c.path, c.transform, c.err, c.errs, got, c.explain)
		}
	}
}

func TestJSONEncodeTransformation(t *testing.T) {
	path := "s:example.net/foo_137x0:737x450_800x600_jpg.webp"
	reference := "../explain_example.json"

	content, err := ioutil.ReadFile(reference)

	if err != nil {
		panic(err)
	}

	want := string(content)

	transform, errs, err := Decode(path, "foo")

	actual := jsonEncodeTransformation("/"+path, transform, errs, err)

	if actual != want {
		t.Errorf("Expected JSON for %v doesn't match with the result saved as %v", path, reference)
	}
}

func TestServerExplain(t *testing.T) {
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

func TestServerRequestBodyPathExplain(t *testing.T) {
	for _, c := range CreateRequestPathCases {
		req, _ := http.NewRequest("GET", "/?explain", bytes.NewBufferString(c.doc))
		w := httptest.NewRecorder()
		http.HandlerFunc(Handler).ServeHTTP(w, req)

		var pi publicImage

		err := json.Unmarshal(w.Body.Bytes(), &pi)

		if pi.Path != c.path || err != nil {
			t.Errorf("Expected path returned by ?explain are invalid or error happened, got %v, want %v (%v)", "blob", pi.Path, c.path, err)
		}
	}
}

func TestServerRequestBodyExplain(t *testing.T) {
	url := "/"
	reference := "../explain_example.json"
	body := `{
		"backend": "https://example.net",
		"path": "foo.jpg",
		"crop": {
			"x": 137,
			"y": 0,
			"width": 737,
			"height": 450
		},
		"width": 800,
		"height": 600,
		"output": "webp"
	}`

	req, _ := http.NewRequest("GET", url+"?explain", bytes.NewBufferString(body))
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
	reference := "../explain_example_single_be.json"

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
	defaultBackend := Backend
	Backend = "https://localhost/"

	for _, c := range BadRequestsURLCases {
		req, _ := http.NewRequest("GET", c.request, nil)

		w := httptest.NewRecorder()
		http.HandlerFunc(Handler).ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusBadRequest)
		}

		if w.Body.String() != BadRequestMessage+"\n" {
			t.Errorf("Bad request body message response is %v, want %v", w.Body.String(), BadRequestMessage)
		}
	}

	Backend = defaultBackend
}

func TestServerBadBodyRequest(t *testing.T) {
	// don't run in parallel due to mocking Backend
	defaultBackend := Backend
	Backend = "https://localhost/"

	for _, c := range BadRequestsURLCases {
		req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(c.request))
		w := httptest.NewRecorder()
		http.HandlerFunc(Handler).ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusBadRequest)
		}

		if w.Body.String() != BadRequestMessage+"\n" {
			t.Errorf("Bad request body message response is %v, want %v", w.Body.String(), BadRequestMessage)
		}
	}

	Backend = defaultBackend
}

func TestServerNotFound(t *testing.T) {
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
			t.Errorf("Error identifying image: want %v, but it was not found for %+v, found instead %+v", v, transform, info)
		}
	}
}

func compareImage(output string, transform image.Transform, width int, height int, t *testing.T) {
	reference := "../test_assets/compare" + transform.Path

	compare := exec.Command("compare", "-metric", "AE", "-fuzz", "15%", reference, output, "/dev/null")

	// always ignore compare exit code
	out, _ := compare.CombinedOutput()
	content := string(out)

	if strings.Index(content, "not supported") != -1 || strings.Index(content, "no decode") != -1 {
		fmt.Println(fmt.Sprintf("Skipping comparing for %v: no support / decoder for ImageMagick compare tool", reference))
		return
	}

	r, _ := regexp.Compile("(?m)^([0-9]+)$")

	abs, err := strconv.Atoi(r.FindString(content))

	if err != nil {
		fmt.Println("compare content:", content)
		panic(err)
	}

	diff := float64(abs) / float64(width*height)

	if diff > 0.15 {
		t.Errorf("Image %+v is very different than expected", transform)
	}
}

func validateGoodRequest(compBackend string, c GoodRequestProvider,
	t *testing.T, w *httptest.ResponseRecorder, defaultOutput string) {
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

func verifyGoodRequestByPath(compBackend string, c GoodRequestProvider, t *testing.T) {
	url := "/" + compBackend + c.url
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	defaultOutput := "jpg"

	if c.acceptWebp {
		req.Header.Set("Accept", "image/webp,*/*;q=0.8")
		defaultOutput = "webp"
	}

	http.HandlerFunc(Handler).ServeHTTP(w, req)
	validateGoodRequest(compBackend, c, t, w, defaultOutput)
}

func verifyGoodRequestByRequestBody(compBackend string, c GoodRequestProvider, t *testing.T) {
	body := strings.Replace(c.requestBody, "REPLACE_ON_TEST", compBackend, -1)
	req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	defaultOutput := "jpg"

	if c.acceptWebp {
		req.Header.Set("Accept", "image/webp,*/*;q=0.8")
		defaultOutput = "webp"
	}

	http.HandlerFunc(Handler).ServeHTTP(w, req)
	validateGoodRequest(compBackend, c, t, w, defaultOutput)
}

func verifyGoodRequest(compBackend string, c GoodRequestProvider, t *testing.T) {
	verifyGoodRequestByPath(compBackend, c, t)
	verifyGoodRequestByRequestBody(compBackend, c, t)
}

func TestServerGoodRequests(t *testing.T) {
	compBackend := compressHost(ts.URL)

	for _, c := range GoodRequestsCases {
		verifyGoodRequest(compBackend, c, t)
	}
}

func TestServerProcessingFailure(t *testing.T) {
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
		got := encodeDimension(string(in.Width), string(in.Height))

		if got != c.want {
			t.Errorf("EncodeDimension(%v, %v) == %v, want %v", in.Width, in.Height, got, c.want)
		}
	}
}

func TestCreateRequestPath(t *testing.T) {
	for _, c := range CreateRequestPathCases {
		path, err := createRequestPath(bytes.NewBufferString(c.doc))

		if path != c.path || err != nil {
			t.Errorf("createRequestPath(%v) == %v, %v want %v %v", "blob", path, err, c.path, nil)
		}
	}
}

func benchmarkGoodRequest(compBackend string, c GoodRequestProvider, t *testing.B) {
	url := "/" + compBackend + c.url

	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	if c.acceptWebp {
		req.Header.Set("Accept", "image/webp,*/*;q=0.8")
	}

	http.HandlerFunc(Handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusOK)
	}

	actualContentType := w.Header().Get("Content-Type")

	if actualContentType != c.outputFiletype {
		t.Errorf("Content-Type is %v, want %v", actualContentType, c.outputFiletype)
	}
}

func BenchmarkServerGoodRequestsPerformance(t *testing.B) {
	compBackend := compressHost(ts.URL)

	defaultBackend := Backend
	Backend = ""

	for _, c := range GoodRequestsCases {
		benchmarkGoodRequest(compBackend, c, t)
	}

	Backend = defaultBackend
}
