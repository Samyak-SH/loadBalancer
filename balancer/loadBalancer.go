package balancer

import (
	"encoding/json"
	"fmt"
	"making-loadbalancer/reqQueue"
	"making-loadbalancer/worker"
	"os"
)

type LoadBalancer struct {
	PORT         int16 `json:"PORT"`
	RequestQueue reqQueue.ReqQueue
	ThreadLimit  int64    `json:"ThreadLimit"`
	Servers      []string `json:"Servers"`
}

func InitializeBalancer(configFilePath string) (*LoadBalancer, error) {
	queueBufferSize := 100
	data, fileReadErr := os.ReadFile(configFilePath)
	if fileReadErr != nil {
		return nil, fileReadErr
	}
	fmt.Println("config file reading complete")
	lb := new(LoadBalancer)
	err := json.Unmarshal(data, lb)

	lb.RequestQueue = reqQueue.NewReqQueue(queueBufferSize)

	if err != nil {
		return nil, err
	}

	fmt.Println("parsing complete")

	return lb, nil

}

func (lb *LoadBalancer) StartWorkers() {
	for i := 0; i < int(lb.ThreadLimit); i++ {
		go func(workerID int) {
			for {
				req := lb.RequestQueue.Dequeue()
				worker.ProcessRequest(workerID, req)
			}
		}(i + 1) // Pass workerID
	}
}
