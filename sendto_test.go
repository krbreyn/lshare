package sendto

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileServer(t *testing.T) {
	var (
		testFile = fileToServe{"helloworld.txt", []byte("Hello, World!")}
	)

	t.Run("all requests should be 404'd if no endpoints have been registered", func(t *testing.T) {
		path := fmt.Sprintf("/file/%s", strings.TrimSuffix(testFile.name, filepath.Ext(testFile.name)))
		request := httptest.NewRequest(http.MethodGet, path, nil)
		responseRecorder := httptest.NewRecorder()

		FileDownloadHandler(responseRecorder, request)

		resp := responseRecorder.Result()
		AssertStatus404(t, resp)
	})

	t.Run("registering the server with a page name, byte data and a file name", func(t *testing.T) {

	})

	t.Run("registering the server with a page name and a file path", func(t *testing.T) {

	})
}

type fileToServe struct {
	name    string
	content []byte
}

func AssertStatusOK(t *testing.T, resp *http.Response) {
	t.Helper()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.Status)
	}
}

func AssertStatus404(t *testing.T, resp *http.Response) {
	t.Helper()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected not found, got %v", resp.Status)
	}
}

func AssertFilenameEqual(t *testing.T, got *http.Response, want string) {
	t.Helper()
	contentDisposition := got.Header.Get("Content-Disposition")
	expectedFilename := fmt.Sprintf("attachment; filename=%s", want)
	if contentDisposition != expectedFilename {
		t.Errorf("expected filename %s, got %s", want, got)
	}
}
