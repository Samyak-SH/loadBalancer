package reqQueue

import (
	"fmt"
	"net/http"
)

type Queue struct {
	Request  *http.Request
	Response http.ResponseWriter
	Done     chan string
}

type ReqQueue struct {
	queue chan Queue
}

func (rq *ReqQueue) Enqueue(r *http.Request, w http.ResponseWriter) {
	fmt.Println("Enquing..")
	done := make(chan string)
	newRequst := Queue{
		Request:  r,
		Response: w,
		Done:     done,
	}
	rq.queue <- newRequst
	fmt.Println("enqueued")
}

func (rq *ReqQueue) Dequeue() Queue {
	return <-rq.queue
}

func NewReqQueue(queueBufferSize int) ReqQueue {
	return ReqQueue{
		queue: make(chan Queue, queueBufferSize),
	}
}
