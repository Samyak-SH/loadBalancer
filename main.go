package main

import (
	"fmt"
	"log"
	loadbalancer "making-loadbalancer/loadBalancer"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
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
		fmt.Printf("\nAlgorithm: Round robin\n\n")
		break
	case 2:
		fmt.Printf("\nAlgorithm: Sticky Session\n\n")
		break
	case 3:
		fmt.Printf("\nAlgorithm: IP Hashing\n\n")
		break
	}

	//handle redirecting
	http.HandleFunc("/", lb.Serve)
	fmt.Printf("Load balancer started on http://localhost:%d\n", lb.PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", lb.PORT), nil))
}
