package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)
// a string
type String string
// a server
func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s)
}
// a greeting
type Struct struct {
	Greeting string
	Punct    string
	Who      string
}
// a server
func (s *Struct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err == nil {
		fmt.Fprint(w, fmt.Sprintf("%s", b))
	} else {
		fmt.Fprint(w, fmt.Sprintf("gads: %s", err))
	}
}

func main() {
	http.Handle("/string", String("I'm a frayed knot."))
	http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
	http.ListenAndServe("localhost:4000", nil)
}
