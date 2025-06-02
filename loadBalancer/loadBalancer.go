package loadbalancer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"making-loadbalancer/server"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LoadBalancer struct {
	PORT                uint16
	Servers             []*server.Server
	Algorithm           uint16
	Proxy               httputil.ReverseProxy
	CurrentServerIndex  int
	ServerCount         int
	HealthCheckInterval int
	SecretKey           string
}

type configFile struct {
	Port                uint16   `json:"PORT"`
	Servers             []string `json:"Servers"`
	Algorithm           uint16   `json:"Algorithm"`
	HealthCheckInterval int      `json:"HealthCheckInterval"`
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
	lb.HealthCheckInterval = cf.HealthCheckInterval
	lb.SecretKey = os.Getenv("SECRET_KEY")
	for _, url := range cf.Servers {
		s := server.NewServer(url)
		lb.Servers = append(lb.Servers, s)
	}

	return lb, nil
}

// Helper fucntions
func (lb *LoadBalancer) getNextServer() (*server.Server, int, error) {
	attempt := 0
	for !lb.Servers[lb.CurrentServerIndex].IsAlive() && attempt < lb.ServerCount {
		lb.CurrentServerIndex = (lb.CurrentServerIndex + 1) % lb.ServerCount
		attempt++
	}
	if attempt > lb.ServerCount {
		return &server.Server{}, -1, errors.New("no healthy servers available")
	}
	server := lb.Servers[lb.CurrentServerIndex]
	assignedServerIndex := lb.CurrentServerIndex
	lb.CurrentServerIndex = (lb.CurrentServerIndex + 1) % lb.ServerCount
	return server, assignedServerIndex, nil
}

func createSignature(value, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	signature := h.Sum(nil)

	signatureHex := hex.EncodeToString(signature)
	return signatureHex
}

func verifySignature(value, key string) (bool, int) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return false, -1
	}
	indexString := parts[0]
	signatureHex := parts[1]

	expectedSignatureHex := createSignature(indexString, key)
	// fmt.Println("Recieved signature", signatureHex)
	// fmt.Println("Expected signature", expectedSignatureHex)
	if !hmac.Equal([]byte(expectedSignatureHex), []byte(signatureHex)) {
		return false, -1
	}
	indexInt, conversionErr := strconv.Atoi(indexString)
	if conversionErr != nil {
		return false, -1
	}
	return true, indexInt
}

// Routing algorithms
func (lb *LoadBalancer) roundRobin(w http.ResponseWriter, r *http.Request) {
	nextServer, _, err := lb.getNextServer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Forwarding request to %s\n", nextServer.GetServerURL())
	nextServer.Serve(w, r)
}

func (lb *LoadBalancer) stickySession(w http.ResponseWriter, r *http.Request) {
	ssidCookie, err := r.Cookie("SSID")
	if err != nil {
		if err == http.ErrNoCookie {
			// fmt.Println("noo cookie ssid found")
			//get server to redirect to
			nextServer, serverIndex, getServerError := lb.getNextServer()
			if getServerError != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//encrypt server index and store it in client's cookie
			serverSignature := createSignature(strconv.Itoa(serverIndex), lb.SecretKey)
			newSsidCookieValue := strconv.Itoa(serverIndex) + "." + serverSignature
			// fmt.Println("newwSsidCOokieValue", newSsidCookieValue)
			newSSIDCookie := &http.Cookie{
				Name:     "SSID",
				Value:    newSsidCookieValue,
				HttpOnly: true,
			}
			http.SetCookie(w, newSSIDCookie)
			//serve user with that server
			nextServer.Serve(w, r)
			return
		} else {
			http.Error(w, "Failed to process cookie", http.StatusBadRequest)
		}
	}
	//Verify signature and serve
	verified, serverIndex := verifySignature(ssidCookie.Value, lb.SecretKey)
	if verified {
		if serverIndex > lb.ServerCount-1 {
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
		} else {
			lb.Servers[serverIndex].Serve(w, r)
		}
	} else {
		http.Error(w, "Invalid cookie", http.StatusBadRequest)
	}
}

// Alogrithm detection
func (lb *LoadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.URL.Path)
	switch lb.Algorithm {
	case 1:
		lb.roundRobin(w, r)
		break
	case 2:
		lb.stickySession(w, r)
		break
	default:
		http.Error(w, "Invalid load balancing algorithm", http.StatusBadRequest)
	}
}

// Health check
func (lb *LoadBalancer) StartHealthChecks(wg *sync.WaitGroup) {
	client := &http.Client{Timeout: 2 * time.Second}
	for i := range lb.Servers {
		go lb.Servers[i].StartHealthCheck(client, wg, lb.HealthCheckInterval)
	}
}
