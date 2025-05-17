package main

import (
	"fmt"
	"log"
	"making-loadbalancer/balancer"
	"net/http"
)

func main() {
	fmt.Println("Started..")
	lb, err := balancer.InitializeBalancer("./config.json")
	if err != nil {
		log.Fatal("failed to initialize loadBalancer, reason\n", err)
	}

	lb.StartWorkers()
	fmt.Printf("load balancer started on http://localhost:%d\n", lb.PORT)
	http.HandleFunc("/", lb.ReqHandler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", lb.PORT), nil))

}
