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
Work performs the work simulation and returns the result in ProcessingResponse.
*/
func Work(pr ProcessingResponse) ProcessingResponse {
	rand.Seed(time.Now().UnixNano())
	durations := []int{5, 50, 500}
	d := durations[rand.Intn(len(durations))]
	time.Sleep(time.Duration(d) * time.Millisecond)
	uuid := uuid.NewRandom()
	log.Printf("worker assigned ResponseId: %s", uuid)
	var resp = ProcessingResponse{
		ResponseId:      uuid.String(),
		ResponseContent: fmt.Sprintf("slept for %dms", d),
		ResponseState:   "in progress",
	}
	return resp
}
