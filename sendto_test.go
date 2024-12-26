package sendto

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileServer(t *testing.T) {
	var (
		testFile = fileToServe{"/helloworld", "helloworld.txt", []byte("Hello, World!")}
		server   = NewFileServer()
	)

	t.Run("all requests should be 404'd if no endpoints have been registered", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, testFile.url, nil)
		responseRecorder := httptest.NewRecorder()

		server.FileDownloadHandler(responseRecorder, request)

		resp := responseRecorder.Result()
		AssertStatus404(t, resp)
	})

	t.Run("registering an endpoint with content", func(t *testing.T) {
		server.RegisterEndpoint(testFile.url, testFile.filename, testFile.content)

		request := httptest.NewRequest(http.MethodGet, testFile.url, nil)
		responseRecorder := httptest.NewRecorder()

		server.FileDownloadHandler(responseRecorder, request)

		resp := responseRecorder.Result()
		AssertStatusOK(t, resp)
		AssertFilenameEqual(t, resp, testFile.filename)

		got := responseRecorder.Body.String()
		want := string(testFile.content)
		if got != want {
			t.Errorf("expected body %s, got %s", want, got)
		}
	})

}

type fileToServe struct {
	url, filename string
	content       []byte
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
		t.Errorf("expected filename %s, got %s", want, contentDisposition)
	}
}
