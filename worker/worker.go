/*
Package worker implements a simulation of a task in a series
of task steps.

It delays, by sleeping, a random duration. It provides a
ProcessingResponse to its invoker.
*/
package worker

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// TODO: this is just to silence compiler's complaints that there's no logging below:
var _ = log.Printf

/*
ProcessingRequest tells the worker what to do. The client or service
that prepares the ProcessingRequest provides the ClientSuppliedRequestId,
normally some kind of UUID.  The service that dispatches the
ProcessingRequest to the worker provides the ServiceSuppliedRequestId,
normally some kind of UUID.  The client or service supplies the
RequestContent, an arbitrary string with meaning to both the client
and the worker.  The service supplies the RequestState.
*/
type ProcessingRequest struct {
	ClientSuppliedRequestId  string
	ServiceSuppliedRequestId string
	RequestContent           string
	RequestState             string
}

/*
Worker uses ProcessingResponse to describe the result of the work.
ResponseId is a UUID.  ResponseContent is an arbitrary string.
ResponseState is a descriptive string.
*/
type ProcessingResponse struct {
	ResponseId      string
	ResponseContent string
	ResponseState   string
}

/*
A client or service prepares a ProcessingOrder to manage a set of
work to be performed to fulfill the request. There is one
ProcessingRequest and a slice of ProcessingResponse elements.
*/
type ProcessingOrder struct {
	Request  ProcessingRequest
	Response []ProcessingResponse
}

/*
Step is the work that worker needs to do. worker.Work waits for a
Step to be added to its queue. then it does the work.
*/
type Step struct {
	Request    ProcessingRequest
	Task       func(ProcessingRequest) ProcessingResponse
	ResultChan chan ProcessingResponse
}

var durations = []time.Duration{
	5 * time.Millisecond,
	50 * time.Millisecond,
	500 * time.Millisecond,
}

/*
Work performs the work simulation and returns the result in ProcessingResponse.
*/
func Work(req ProcessingRequest) ProcessingResponse {
	rand.Seed(time.Now().UnixNano())
	d := durations[rand.Intn(len(durations))]
	time.Sleep(d)
	uuid := uuid.NewRandom()
	return ProcessingResponse{
		ResponseId:      uuid.String(),
		ResponseContent: fmt.Sprintf("slept for %dms", d/time.Millisecond),
		ResponseState:   "done",
	}
}

/*
handle() dispatches the queued work.
*/
func Handle(queue chan *Step) {
	for req := range queue {
		req.ResultChan <- req.Task(req.Request)
		close(req.ResultChan)
	}
}
