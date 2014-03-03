package main

import (
	"encoding/json"
	"fmt"
	"github.com/musicbeat/worker"
	"log"
	"net/http"
)

type ClientRequest struct {
	RequestId      string
	RequestContent string
	RequestState   string
}

func (req ClientRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// prepare the processing order:
	po := worker.ProcessingOrder{
		Request: worker.ProcessingRequest{
			ClientSuppliedRequestId: req.RequestId,
			RequestContent:          req.RequestContent,
			RequestState:            req.RequestState,
		},
		Response: []worker.ProcessingResponse{},
	}
	c := make(chan int)
	// just execute a work function:
	go worker.Work(&po, c)
	// wait for it:
	<- c
	j, err := json.MarshalIndent(po, "", "  ")
	if err == nil {
		fmt.Fprint(w, fmt.Sprintf("%s\n", j))
	} else {
		fmt.Fprint(w, fmt.Sprintf("gads: %s\n", err))
	}
}

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
