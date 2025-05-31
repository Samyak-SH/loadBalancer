package main

import (
	"fmt"
	"log"
	loadbalancer "making-loadbalancer/loadBalancer"
	"net/http"
)

func main() {
	lb, err := loadbalancer.Initialize("./config.json")
	if err != nil {
		log.Fatal("failed to initialize load balancer, reason\n", err)
	}

	http.HandleFunc("/", lb.Serve)
	fmt.Printf("Load balancer started on http://localhost:%d\n", lb.PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", lb.PORT), nil))
}
