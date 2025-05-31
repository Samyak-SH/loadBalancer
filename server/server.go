package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server struct {
	serverURL string
	isAlive   bool
	proxy     *httputil.ReverseProxy
}

func NewServer(address string) *Server {
	serverUrl, err := url.Parse(address)
	if err != nil {
		log.Fatal("failed to parse address of a server ", address)
		os.Exit(0)
	}
	return &Server{
		serverURL: address,
		isAlive:   true,
		proxy:     httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func (s *Server) IsAlive() bool {
	return s.isAlive
}

func (s *Server) GetServerURL() string {
	return s.serverURL
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}
