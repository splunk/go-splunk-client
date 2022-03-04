package main

import (
	"fmt"
	"reflect"
)

type myStruct struct {
	priv string
	Pub  string
	Nest struct {
		Again struct {
		}
	}
}

func main() {
	origPtr := &myStruct{}

	origTPtr := reflect.TypeOf(origPtr)
	origT := origTPtr.Elem()
	origFields := make([]reflect.StructField, origT.NumField())
	for i := 0; i < origT.NumField(); i++ {
		origFields[i] = origT.Field(i)
	}
	newT := reflect.StructOf(origFields)
	newVPtr := reflect.New(newT)
	newV := newVPtr.Elem()
	fmt.Printf("new: %#v\n", newV)

	origVPtr := reflect.ValueOf(origPtr)
	origV := origVPtr.Elem()
	origV.Set(newV)
	fmt.Printf("orig: %#v\n", *origPtr)
}
