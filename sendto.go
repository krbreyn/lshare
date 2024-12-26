package sendto

import (
	"log"
	"net/http"
	"sync"
)

type FileServer struct {
	mu        sync.RWMutex
	endpoints map[string]endpoint
}

type endpoint struct {
	url      string
	filename string
	content  []byte
}

func NewFileServer() *FileServer {
	return &FileServer{
		endpoints: make(map[string]endpoint),
	}
}

func (s *FileServer) StartServer(port string) {
	http.HandleFunc("/", s.FileDownloadHandler)
	if port == "" {
		port = ":8080"
	}

	log.Printf("starting server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func (s *FileServer) RegisterEndpoint(url, filepath string, content []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.endpoints[url] = endpoint{url, filepath, content}
}

func (s *FileServer) DeleteEndpoint(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.endpoints, url)
}

func (s *FileServer) FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ep, ok := s.endpoints[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+ep.filename)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(ep.content); err != nil {
		return
	}
}
