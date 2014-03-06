package main

import (
	"encoding/json"
	"fmt"
	"github.com/musicbeat/go/worker"
	"log"
	"net/http"
	"time"
)

type ClientRequest struct {
	RequestId      string
	RequestContent string
	RequestState   string
}

const (
	stepCount  = 3
	sleepCount = 10
	sleepTime  = time.Millisecond * 10
)

/*
ServeHTTP represents the processing service. It implements the
series of business logic required to service the ClientRequest.
The business logic consists of performing the same step stepCount
times.  Each time the step is performed it has a randomly selected
latency.  The business logic steps are placed on a queue. worker.Work()
performs them in succession. ServeHTTP sleeps for 10ms, checks to
see if everything is done. If not, it sleeps for another 10ms. If,
after 10 iterations of 10ms sleeping (100ms or so), it declares
itself done enough.
*/
func (req ClientRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	// prepare the processing order:
	po := worker.ProcessingOrder{
		Request: worker.ProcessingRequest{
			ClientSuppliedRequestId: req.RequestId,
			RequestContent:          req.RequestContent,
			RequestState:            req.RequestState,
		},
		Response: []worker.ProcessingResponse{},
	}

	queue := make(chan *worker.Step, stepCount)

	// set up the multi-step business process:
	var bp []worker.Step = make([]worker.Step, stepCount)
	for i := 0; i < stepCount; i++ {
		bp[i] = worker.Step{po.Request, worker.Work, make(chan worker.ProcessingResponse)}
		queue <- &bp[i]
	}
	// put the worker to work:
	go worker.Handle(queue)

	// collect the responses, but don't wait too long:
	c := make(chan uint)
	go collect(bp, &po, c, start)
	bits := <- c
	if bits == 0 {
		// note that the request is completed:
		po.Request.RequestState = "complete"
	} else {
		po.Request.RequestState = "in progress"
	}

	log.Printf("elapsed: %03d\n", time.Since(start).Nanoseconds()/1e6)

	// jsonify:
	j, err := json.MarshalIndent(po, "", "  ")
	if err == nil {
		fmt.Fprint(w, fmt.Sprintf("%s\n", j))
	} else {
		fmt.Fprint(w, fmt.Sprintf("gads: %s\n", err))
	}
}

func collect(bp []worker.Step, po *worker.ProcessingOrder, c chan uint, start time.Time) {
	// keep track in a bitmap:
	var bits, i uint
	for i = 0; i < stepCount; i++ {
		bits = bits | 1 << i
	}
	for i := 0; i < sleepCount; i++ {
		// sleep
		<-time.NewTimer(sleepTime).C
		// log.Printf("collect elapsed: %03d\n", time.Since(start).Nanoseconds()/1e6)
		// check each step in the business process:
		var j uint
		for j = 0; j < stepCount; j++ {
			// do a non-blocking read on the channel
			select {
			case pr, found := <-bp[j].ResultChan:
				// if there's a response, append it to the processing order's responses
				if found {
					po.Response = append(po.Response, pr)
					// clear this step's bit:
					bits = bits ^ 1 << j
					log.Printf("poll %d; chan %d; bits: %b\n", i, j, bits)
				}
			default:
				// log.Printf("poll %d; chan %d; nothing yet\n", i, j)
			}
		}
		if bits == 0 {
			break
		}
	}
	// log.Printf("collect elapsed: %03d\n", time.Since(start).Nanoseconds()/1e6)
	c <- bits
}

/*
main simulates the client. It looks like it's waiting for a
client to present a request. But for this simulation, it
fakes the ClientRequest itself. It waits for a client to do
a get on the /process resource.
*/
func main() {
	r := ClientRequest{
		RequestId:      "1001",
		RequestContent: "all your gifts are belong to us",
		RequestState:   "submitted",
	}
	http.Handle("/process", &r)
	err := http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
