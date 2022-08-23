package main

import (
	"fmt"
	"io"
	"net/http"
)

func Greet(writer io.Writer, name string) {
	_, err := fmt.Fprintf(writer, "<html><body><h>Hello, %s</h></body></html>", name)
	if err != nil {
		return
	}
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "world")
}

//func main() {
//	log.Fatal(http.ListenAndServe(":5001", http.HandlerFunc(MyGreeterHandler)))
//
//}
