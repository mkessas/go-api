package main

import "net/http"

// Path stores a router path and associated handler
type Path struct {
	path      string
	handler   func(http.ResponseWriter, *http.Request)
	method    string
	protected bool
}

// Creds stores login credentials
type Creds struct {
	username string
	password string
}
