package main

import (
	"net/http"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

//provide context for each request.
type Context struct {
	*http.Request
	http.ResponseWriter
}

func main() {
	println(`Go's field is backend.In this exercise,we focus on web framework,it shortens our development time and reduces our coding work.
Currently,there are 2 mainstream implementations,function handler and struct handler.
We will implement a function handler tiny webframework.
The first part is context management.Next part will focus on middleware`)
}
