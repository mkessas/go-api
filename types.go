package main

import "net/http"

type path struct {
	path      string
	handler   func(http.ResponseWriter, *http.Request)
	method    string
	protected bool
}
