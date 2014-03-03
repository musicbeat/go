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
	"math/rand"
	"time"
)

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

var durations = []int{5, 50, 500}

/*
Work performs the work simulation and returns the result in ProcessingResponse.
*/
func Work(po *ProcessingOrder, c chan int) {
	rand.Seed(time.Now().UnixNano())
	d := durations[rand.Intn(len(durations))]
	time.Sleep(time.Duration(d) * time.Millisecond)
	uuid := uuid.NewRandom()
	resp := ProcessingResponse{
		ResponseId:      uuid.String(),
		ResponseContent: fmt.Sprintf("slept for %dms", d),
		ResponseState:   "done",
	}
	po.Response = append(po.Response, resp)
	c <- 1
	return
}
