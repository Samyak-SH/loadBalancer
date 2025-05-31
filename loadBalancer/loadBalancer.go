package loadbalancer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"making-loadbalancer/server"
	"net/http"
	"net/http/httputil"
	"os"
)

type LoadBalancer struct {
	PORT               uint16
	Servers            []server.Server
	Algorithm          uint16
	Proxy              httputil.ReverseProxy
	CurrentServerIndex int
	ServerCount        int
}

type configFile struct {
	Port      uint16   `json:"PORT"`
	Servers   []string `json:"Servers"`
	Algorithm uint16   `json:"Algorithm"`
}

func Initialize(configFilePath string) (*LoadBalancer, error) {

	//read config.json file and store marshalled data in "data"
	data, fileReadError := os.ReadFile((configFilePath))
	if fileReadError != nil {
		return nil, fileReadError
	}

	//unmarshal json data from "data" into cf struct
	cf := new(configFile)
	parsingError := json.Unmarshal(data, cf)
	if parsingError != nil {
		return nil, parsingError
	}

	//intialize loadBalancer
	lb := new(LoadBalancer)
	lb.PORT = cf.Port
	lb.Algorithm = cf.Algorithm
	lb.CurrentServerIndex = 0
	lb.ServerCount = len(cf.Servers)
	for _, url := range cf.Servers {
		s := server.NewServer(url)
		lb.Servers = append(lb.Servers, *s)
	}

	return lb, nil
}

func (lb *LoadBalancer) getNextServer() (server.Server, error) {
	attempt := 0
	for !lb.Servers[lb.CurrentServerIndex].IsAlive() && attempt <= lb.ServerCount {
		lb.CurrentServerIndex = (lb.CurrentServerIndex + 1) % lb.ServerCount
		attempt++
	}
	if attempt > lb.ServerCount {
		return server.Server{}, errors.New("no healthy servers available")
	}
	server := lb.Servers[lb.CurrentServerIndex]
	lb.CurrentServerIndex = (lb.CurrentServerIndex + 1) % lb.ServerCount
	return server, nil
}

func (lb *LoadBalancer) roundRobin(w http.ResponseWriter, r *http.Request) {
	nextServer, err := lb.getNextServer()
	fmt.Println("Serving using ", nextServer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Forwarding request to %s\n", nextServer.GetServerURL())
	nextServer.Serve(w, r)
}

func (lb *LoadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	switch lb.Algorithm {
	case 1:
		lb.roundRobin(w, r)
		break
	default:
		http.Error(w, "Invalid load balancing algorithm", http.StatusBadRequest)
	}
}
