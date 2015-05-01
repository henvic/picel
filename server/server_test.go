package picel

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

type ExistsDependencyProvider struct {
	cmd  string
	find bool
}

type CheckMissingDependenciesProvider struct {
	cmds      []string
	allExists bool
}

type GoodRequestProvider struct {
	url            string
	outputFiletype string
	meta           []string
	acceptWebp     bool
}

type BuildExplainProvider struct {
	transform Transform
	err       error
	errs      []error
	explain   Explain
}

type ServerProcessingFailureProvider struct {
	url string
}

func init() {
	// binary test assets are stored in a helper branch for neatness
	exec.Command("git", "checkout", "test_assets", "--", "test_assets").Run()
	exec.Command("git", "rm", "--cached", "-r", "test_assets").Run()
}

func TestVersion(t *testing.T) {
	flagVersion = true
	main()
}

func TestExistsDependency(t *testing.T) {
	t.Parallel()
	for _, c := range existsDependencyCases {
		exists := existsDependency(c.cmd)

		if exists != c.find {
			t.Errorf("existsDependency(%v) should return %v", c.cmd, c.find)
		}
	}
}

func TestCheckMissingDependencies(t *testing.T) {
	t.Parallel()
	for _, c := range CheckMissingDependencies {
		var StdoutMock bytes.Buffer
		var StderrMock bytes.Buffer

		defaultStdout, defaultStderr := std.out, std.err
		std.out = log.New(&StdoutMock, "", log.LstdFlags)
		std.err = log.New(&StderrMock, "", log.LstdFlags)
		checkMissingDependencies(c.cmds...)
		std.out, std.err = defaultStdout, defaultStderr

		if StdoutMock.String() != "" {
			t.Errorf("checkMissingDependencies(%v) stdout should be empty", c.cmds)
		}

		// errMessage := StderrMock.String()

		// if
		// "Dependencies missing:", strings.Join(missing, ", ")

		// fmt.Println(StderrMock.String())

		// if !existsDependencies(c.cmds...) {
		// 	t.Errorf("existsDependencies(%v) should return true", c.cmds)
		// }

	}
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
	url := "foo_137x0:737x450_800x600_jpg.webp"
	reference := "explain_example.json"

	content, err := ioutil.ReadFile(reference)

	if err != nil {
		panic(err)
	}

	want := string(content)

	actual := jsonEncodeTransformation(Decode(url, "foo"))
	jsonEncodeTransformation(Transform{}, nil, nil)

	if actual != want {
		t.Errorf("Expected JSON for %v doesn't match with the result saved as %v", url, reference)
	}
}

func TestServerExplain(t *testing.T) {
	t.Parallel()
	url := "/foo_137x0:737x450_800x600_jpg.webp"
	reference := "explain_example.json"

	req, _ := http.NewRequest("GET", url+"?explain", nil)
	w := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(w, req)

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
	t.Parallel()
	url := "/_"

	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(w, req)

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
	http.HandlerFunc(handler).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusNotFound)
	}
}

func verifyGoodRequest(c GoodRequestProvider, t *testing.T) {
	req, _ := http.NewRequest("GET", c.url, nil)
	w := httptest.NewRecorder()
	defaultOutput := "jpg"

	if c.acceptWebp {
		req.Header.Set("Accept", "image/webp,*/*;q=0.8")
		defaultOutput = "webp"
	}

	http.HandlerFunc(handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusOK)
	}

	actualContentType := w.Header().Get("Content-Type")

	if actualContentType != c.outputFiletype {
		t.Errorf("Content-Type is %v, want %v", actualContentType, c.outputFiletype)
	}

	// this checking is very unsafe as can lead to false positives
	// but works with the current test cases
	transform, _, _ := Decode(c.url, defaultOutput)

	if transform.Raw {
		_, filename := transform.Image.name()
		reference, _ := ioutil.ReadFile("test_assets/" + filename)

		if bytes.Compare(reference, w.Body.Bytes()) != 0 {
			t.Errorf("Raw file for %v differ from what is expected", filename)
		}
	}

	// identify doesn't support some standards like .gif
	if c.meta == nil {
		return
	}

	identify := exec.Command("identify", "-verbose", "-")
	identify.Stdin = w.Body
	out, identifyImageErr := identify.CombinedOutput()
	info := string(out)

	if identifyImageErr != nil {
		t.Errorf("Error while identifying image: %v", identifyImageErr)
	}

	for _, v := range c.meta {
		if strings.LastIndex(info, v) == -1 {
			t.Errorf("Error identifying image: want %v, but it was not found", v)
		}
	}
}

func TestServerGoodRequests(t *testing.T) {
	t.Parallel()
	fsHandler := func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "../", "", -1)
		http.ServeFile(w, r, "test_assets/"+path[1:])
	}

	ts := httptest.NewServer(http.HandlerFunc(fsHandler))
	defer ts.Close()

	backend = ts.URL + "/"

	for _, c := range GoodRequestsCases {
		verifyGoodRequest(c, t)
	}
}

func TestServerProcessingFailure(t *testing.T) {
	t.Parallel()
	fsHandler := func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "../", "", -1)
		http.ServeFile(w, r, "test_assets/"+path[1:])
	}

	ts := httptest.NewServer(http.HandlerFunc(fsHandler))
	defer ts.Close()

	backend = ts.URL + "/"

	for _, c := range ServerProcessingFailureCases {
		req, _ := http.NewRequest("GET", c.url, nil)
		w := httptest.NewRecorder()

		http.HandlerFunc(handler).ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Request status code response is %v, want %v", w.Code, http.StatusInternalServerError)
		}
	}
}
