package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
)

var creds = Creds{username: "admin", password: "admin"}

var paths = []Path{
	{path: "/api/test1", handler: test1Handler, method: "GET", protected: true},
	{path: "/api/test2", handler: test2Handler, method: "GET", protected: false},
	{path: "/api/hello/{name}", handler: helloHandler, method: "GET", protected: false},
	{path: "/api/temperature", handler: temperatureHandler, method: "GET", protected: false},
}

func test1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test 1 handler!")
}

func test2Handler(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(413)
	w.Header().Set("x-Yes", "Cool")
	fmt.Fprintf(w, "Test 2 handler!\n")

	for k, v := range r.Header {
		fmt.Fprintf(w, "- %s : %s\n", k, v)

	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", vars["name"])
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("http://alarm:3000/api/sensor/status")

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}

	defer r.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	temperature, _, _, err := jsonparser.Get(data, "details", "status", "temperature")

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
	}

	fmt.Fprintf(w, "%s", temperature)

}
