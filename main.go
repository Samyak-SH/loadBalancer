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
	wg.Add(lb.ServerCount)
	lb.StartHealthChecks(&wg)
	wg.Wait()

	http.HandleFunc("/", lb.Serve)
	fmt.Printf("Load balancer started on http://localhost:%d\n", lb.PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", lb.PORT), nil))
}
