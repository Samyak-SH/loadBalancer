package main

import (
	"fmt"
	"log"
	loadbalancer "making-loadbalancer/loadBalancer"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	var wg sync.WaitGroup
	//load env files
	wg.Add(1)
	envLoadErr := godotenv.Load(".env")
	if envLoadErr != nil {
		log.Fatal("failed to load env files, reason\n", envLoadErr)
	}
	fmt.Println("env files loaded successfully")
	wg.Done()
	//initialize load balancer configurations
	lb, err := loadbalancer.Initialize("./config.json")
	if err != nil {
		log.Fatal("failed to initialize load balancer, reason\n", err)
	}

	//intialize health checks
	wg.Add(lb.ServerCount)
	lb.StartHealthChecks(&wg)
	wg.Wait()
	switch lb.Algorithm {
	case 1:
		fmt.Println("Algorithm: Round robin")
		break
	case 2:
		fmt.Println("Algorithm: sticky session")
		break
	}

	//handle redirecting
	http.HandleFunc("/", lb.Serve)
	fmt.Printf("Load balancer started on http://localhost:%d\n", lb.PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", lb.PORT), nil))
}
