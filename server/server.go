package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

type Server struct {
	serverURL string
	isAlive   bool
	proxy     *httputil.ReverseProxy
	mu        sync.RWMutex
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

func (s *Server) StartHealthCheck(client *http.Client, wg *sync.WaitGroup, healthCheckInterval int) {
	log.Printf("Health check started for %s\n", s.GetServerURL())
	wg.Done()
	for {
		response, err := client.Get(s.GetServerURL())
		if err != nil || (response != nil && response.StatusCode >= 500) {
			s.SetAlive(false)
			log.Printf("Server with url %s down\n", s.GetServerURL())
		} else {
			if !s.isAlive {
				s.SetAlive(true)
				log.Printf("Server with url %s back up\n", s.GetServerURL())
			}

		}
		if response != nil {
			response.Body.Close()
		}
		time.Sleep(time.Duration(healthCheckInterval) * time.Second)
	}
}

func (s *Server) IsAlive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isAlive
}

func (s *Server) SetAlive(status bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isAlive = status
}

func (s *Server) GetServerURL() string {
	return s.serverURL
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}
