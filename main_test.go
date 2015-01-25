package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestTiny(t *testing.T) {

	recorder := httptest.NewRecorder()
	//send a request with post form name=foo
	v := url.Values{}
	v.Set("name", "foo")
	req, _ := http.NewRequest("POST", "http://localhost:3000/hello", strings.NewReader(v.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	s := NewServer()

	//use middleware to parse the form before each request.
	s.Use(new(ParseForm))

	s.Post("/hello", func(ctx *Context) string {
		name := ctx.Form.Get("name")
		return name
	})

	s.ServeHTTP(recorder, req)
	name, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(name) != "foo" {
		t.Fail()
		return
	}
}
