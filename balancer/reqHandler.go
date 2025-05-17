package balancer

import (
	"fmt"
	"net/http"
)

func (lb *LoadBalancer) ReqHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling req")
		lb.RequestQueue.Enqueue(r, w)
	}
}
