package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTiny(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:3000/hello?name=foo", nil)
	s := NewServer()
	s.Get("/hello", func(ctx *Context) string {
		name := ctx.URL.Query().Get("name")
		return name
	})
	s.ServeHTTP(recorder, req)
	name, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(name) != "foo" {
		println(string(name))
		t.Fail()
		return
	}
}
