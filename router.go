package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

var writerForKafka *kafka.Writer

func router() {
	r := mux.NewRouter()
	writerForKafka = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "GlobalExtracgtorTopic",
		Balancer: &kafka.LeastBytes{},
	})
	for _, p := range paths {

		fmt.Printf("Adding handler for '%s'\n", p.path)
		//http.HandleFunc(p.path, p.handler)
		if p.protected {
			h := http.HandlerFunc(p.handler)
			a := httpauth.SimpleBasicAuth(creds.username, creds.password)(h)
			r.Handle(p.path, a).Methods(p.method)
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
