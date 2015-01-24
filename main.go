package main

import (
	"reflect"
)

func MakeMap(fpt interface{}) {
	fnV := reflect.ValueOf(fpt).Elem()
	fnI := reflect.MakeFunc(fnV.Type(), implMap)
	fnV.Set(fnI)
}

//TODO:completes implMap function.
var implMap func([]reflect.Value) []reflect.Value

func main() {

	println("It is said that Go has no generics.\nHowever we have many other ways to implement a generics like library if less smoothly,one is reflect.MakeFunc.\nUnderscore is a very useful js library,and now let's implement part of it-map,it will help you to understand how reflect works.\nPlease finish the 'implMap' function and pass the test.")
}
