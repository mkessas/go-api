package main

import "net/http"

type path struct {
	path      string
	handler   func(http.ResponseWriter, *http.Request)
	method    string
	protected bool
}

var paths = []path{
	{path: "/api/test1", handler: test1Handler, method: "GET", protected: true},
	{path: "/api/test2", handler: test2Handler, method: "GET", protected: false},
}
