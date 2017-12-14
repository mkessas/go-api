package main

var creds = Creds{username: "admin", password: "admin"}

var paths = []Path{
	{path: "/api/test1", handler: test1Handler, method: "GET", protected: true},
	{path: "/api/test2", handler: test2Handler, method: "GET", protected: false},
}
