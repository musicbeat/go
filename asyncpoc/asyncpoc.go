package main

import (
	"encoding/json"
	"fmt"
	"github.com/musicbeat/worker"
	"log"
	"net/http"
)

type ProcessingRequest struct {
	ClientSuppliedRequestId  int
	ServiceSuppliedRequestId int
	RequestContent           string
	RequestState             string
}

type ProcessingOrder struct {
	Request  ProcessingRequest
	Response []worker.ProcessingResponse
}

func (req ProcessingRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// prepare the processing order:
	po := ProcessingOrder{
		Request:  req,
		Response: []worker.ProcessingResponse{},
	}
	// just execute a work function:
	resp := worker.Work(worker.ProcessingResponse{})
	po.Response = append(po.Response, resp)
	j, err := json.MarshalIndent(po, "", "  ")
	if err == nil {
		fmt.Fprint(w, fmt.Sprintf("%s\n", j))
		log.Printf("%+v\n", po)
	} else {
		fmt.Fprint(w, fmt.Sprintf("gads: %s\n", err))
		log.Printf("%s\n", err)
	}
}

func main() {
	r := ProcessingRequest{
		ClientSuppliedRequestId: 1001,
		RequestContent:          "all your gifts are belong to us",
		RequestState:            "submitted",
	}
	http.Handle("/process", &r)
	err := http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
