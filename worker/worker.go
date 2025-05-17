package worker

import (
	"fmt"
	"making-loadbalancer/reqQueue"
	"time"
)

func ProcessRequest(workerID int, req reqQueue.Queue) {
	fmt.Printf("Worker %d processing request: %s\n", workerID, req.Request.URL.Path)
	fmt.Printf("\nworker %d sleeping ............\n............\n............\n............", workerID)
	time.Sleep(2 * time.Second)
	fmt.Printf("\nworker %d awake", workerID)
	fmt.Fprintf(req.Response, "req recieved")
}
