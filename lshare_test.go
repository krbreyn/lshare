package lshare

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

		resp, _ := RecordResponse(t, request, server)
		AssertStatus404(t, resp)
	})

	t.Run("registering an endpoint with content serves direct download", func(t *testing.T) {
		server := NewFileServer()
		server.RegisterEndpoint(testFile.url, testFile.filename, testFile.content)

		request := httptest.NewRequest(http.MethodGet, testFile.url, nil)

		resp, body := RecordResponse(t, request, server)
		AssertFileServed(t, testFile, resp, body)
	})

	t.Run("deleting an endpoint returns 404 on further requests", func(t *testing.T) {
		server := NewFileServer()
		server.RegisterEndpoint(testFile.url, testFile.filename, testFile.content)

		request := httptest.NewRequest(http.MethodGet, testFile.url, nil)

		resp, body := RecordResponse(t, request, server)
		AssertFileServed(t, testFile, resp, body)

		server.DeleteEndpoint(testFile.url)

		resp, _ = RecordResponse(t, request, server)
		AssertStatus404(t, resp)
	})
}

type fileToServe struct {
	url, filename string
	content       []byte
}

func RecordResponse(t *testing.T, r *http.Request, s *FileServer) (*http.Response, string) {
	responseRecorder := httptest.NewRecorder()
	s.FileDownloadHandler(responseRecorder, r)
	resp := responseRecorder.Result()
	return resp, responseRecorder.Body.String()
}

func AssertFileServed(t *testing.T, file fileToServe, resp *http.Response, body string) {
	t.Helper()
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
