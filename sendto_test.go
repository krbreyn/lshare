package sendto

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func TestGetLocalIP(t *testing.T) {
// 	want := "" // local ip would go here
// 	got, err := GetLocalIP()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if got != want {
// 		t.Errorf("got %v, and %v", got, want)
// 	}
// }

func TestFileServer(t *testing.T) {
	var (
		testFile = fileToServe{"/helloworld", "helloworld.txt", []byte("Hello, World!")}
	)

	t.Run("all requests should be 404'd if no endpoints have been registered", func(t *testing.T) {
		server := NewFileServer()
		request := httptest.NewRequest(http.MethodGet, testFile.url, nil)
		responseRecorder := httptest.NewRecorder()

		server.FileDownloadHandler(responseRecorder, request)

		resp := responseRecorder.Result()
		AssertStatus404(t, resp)
	})

	t.Run("registering an endpoint with content", func(t *testing.T) {
		server := NewFileServer()
		server.RegisterEndpoint(testFile.url, testFile.filename, testFile.content)

		request := httptest.NewRequest(http.MethodGet, testFile.url, nil)
		responseRecorder := httptest.NewRecorder()

		server.FileDownloadHandler(responseRecorder, request)
		resp := responseRecorder.Result()
		AssertFileServed(t, testFile, resp, responseRecorder.Body.String())
	})

	t.Run("deleting an endpoint", func(t *testing.T) {
		server := NewFileServer()
		server.RegisterEndpoint(testFile.url, testFile.filename, testFile.content)

		request := httptest.NewRequest(http.MethodGet, testFile.url, nil)
		responseRecorder := httptest.NewRecorder()

		server.FileDownloadHandler(responseRecorder, request)
		resp := responseRecorder.Result()
		AssertFileServed(t, testFile, resp, responseRecorder.Body.String())

		server.DeleteEndpoint(testFile.url)
		responseRecorder = httptest.NewRecorder()
		server.FileDownloadHandler(responseRecorder, request)
		resp = responseRecorder.Result()
		AssertStatus404(t, resp)
	})
}

type fileToServe struct {
	url, filename string
	content       []byte
}

func AssertFileServed(t *testing.T, file fileToServe, resp *http.Response, body string) {
	AssertStatusOK(t, resp)
	AssertFilenameEqual(t, resp, file.filename)

	want := string(file.content)
	if body != want {
		t.Errorf("expected body %s, got %s", want, body)
	}

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
