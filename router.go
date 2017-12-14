package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
)

func protect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Protecting request")
		h.ServeHTTP(w, r)
	})
}

func router() {
	r := mux.NewRouter()

	for _, p := range paths {

		fmt.Printf("Adding handler for '%s'\n", p.path)
		//http.HandleFunc(p.path, p.handler)
		if p.protected {
			r.Handle(p.path, httpauth.SimpleBasicAuth("admin", "admin")(http.HandlerFunc(p.handler))).Methods(p.method)
		} else {
			r.HandleFunc(p.path, p.handler).Methods(p.method)
		}
	}

	r.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "File Not Found")
	})
	r.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":"+strconv.Itoa(port), r)
}

func test1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test 1 handler!")
}

func test2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test 2 handler!")
}
