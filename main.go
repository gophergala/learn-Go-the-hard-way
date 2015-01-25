package main

import (
	"log"
	"net"
	"net/http"
	"reflect"
)

type Server struct {
	routes []route      //routes
	addr   string       //address
	l      net.Listener //save the listener so it can be closed.
}

type route struct {
	r           string        //route url
	method      string        //method (GET)
	httpHandler http.Handler  //custome handler is allowed.
	handler     reflect.Value //handle func
}

//implements http.Handle
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, r := range s.routes {
		if r.r == req.URL.Path {
			if r.httpHandler != nil {
				r.httpHandler.ServeHTTP(res, req)
			} else {
				//TODO:pass the contex to the function and write return value to res.
			}
		}
	}
}

//close the server
func (s *Server) Close() {
	if s.l != nil {
		s.l.Close()
	}
}

//run the server
func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.Handle("/", s)
	log.Printf("start serverving...\nPlease visit http://localhost:3000")
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	s.l = l
	err = http.Serve(s.l, mux)
	s.l.Close()
}

//Get adds a handler for the 'GET' http method for server.
func (s *Server) Get(rt string, handler interface{}) {
	switch handler.(type) {
	case http.Handler:
		s.routes = append(s.routes, route{r: rt, method: "GET", httpHandler: handler.(http.Handler)})
	case reflect.Value:
		fv := handler.(reflect.Value)
		s.routes = append(s.routes, route{r: rt, method: "GET", handler: fv})
	default:
		fv := reflect.ValueOf(handler)
		s.routes = append(s.routes, route{r: rt, method: "GET", handler: fv})
	}
}

func NewServer() *Server {
	return &Server{addr: "localhost:3000"}
}

//provide context for each request.
type Context struct {
	*http.Request
	http.ResponseWriter
	Server *Server
	Params map[string]string
}

var contextType reflect.Type

func init() {
	contextType = reflect.TypeOf(Context{})
}

func main() {
	println(`Go's field is backend.In this exercise,we focus on web framework,it shortens our development time and reduces our coding work.
We will implement a function handler based tiny webframework.
The first part is context management.Next part will focus on middleware.
Now edit main.go file to Complete 'Server.ServeHttp()'.
In this method you need to call the handler in the context and pass context as paramater if in the signature,
and also write the return valut to the responseWriter. `)
}
