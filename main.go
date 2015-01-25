package main

import (
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
)

type Server struct {
	middlewares []Middleware //middleware
	routes      []route      //routes
	addr        string       //address
	l           net.Listener //save the listener so it can be closed.
}

type route struct {
	r           string        //route url
	method      string        //method (GET)
	httpHandler http.Handler  //custome handler is allowed.
	handler     reflect.Value //handle func
}

type Middleware struct {
	Handle (*Context)
}

func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

//implements http.Handle
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, r := range s.routes {
		if r.r == req.URL.Path {
			if r.httpHandler != nil {
				r.httpHandler.ServeHTTP(res, req)
			} else {
				//function handler
				//*context must be the first argument.
				ctx := &Context{req, res, s, make(map[string]string)}
				if err := ctx.ParseForm(); err != nil {
					log.Println(err)
				}
				var args []reflect.Value
				if requiresContext(r.handler.Type()) {
					args = append(args, reflect.ValueOf(ctx))
				}
				ret := r.handler.Call(args)
				if len(ret) == 0 {
					return
				}
				//if has return value,write to response.
				sval := ret[0]
				var content []byte
				if sval.Kind() == reflect.String {
					content = []byte(sval.String())
				} else if sval.Kind() == reflect.Slice && sval.Type().Elem().Kind() == reflect.Uint8 {
					content = sval.Interface().([]byte)
				}
				ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
				ctx.ResponseWriter.Write(content)
			}
		}
	}
}

// SetHeader sets a response header. If `unique` is true, the current value
// of that header will be overwritten . If false, it will be appended.
func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.ResponseWriter.Header().Set(hdr, val)
	} else {
		ctx.ResponseWriter.Header().Add(hdr, val)
	}
}

// requiresContext determines whether 'handlerType' contains
// an argument to 'web.Ctx' as its first argument
func requiresContext(handlerType reflect.Type) bool {
	//if the method doesn't take arguments, no
	if handlerType.NumIn() == 0 {
		return false
	}

	//if the first argument is not a pointer, no
	a0 := handlerType.In(0)
	if a0.Kind() != reflect.Ptr {
		return false
	}
	//if the first argument is a context, yes
	if a0.Elem() == contextType {
		return true
	}

	return false
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
	println(`Now we have a context for each request.
To decouple the request handling,middleware is very useful,each middleware just deals with part of the handling,
and pass the control to the next,middleware is pluggable.
That means you can just add middlewares you want to use to handle the request.
Suppose you have done the task l6,and have a context,now we need to add a middleware layer to tiny webframework,and make it pluggable.
Edit main.go,and finish the task.`)
}
