package main

import (
	"fmt"
	"log"
	loadbalancer "making-loadbalancer/loadBalancer"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	lb, err := loadbalancer.Initialize("./config.json")
	if err != nil {
		log.Fatal("failed to initialize load balancer, reason\n", err)
	}
	envs := os.Environ()
	for index, env := range envs {
		splitenv := strings.Split(env, "=")
		fmt.Printf("%d %s\n", index, splitenv[0])
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
