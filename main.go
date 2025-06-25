package main

import (
	"fmt"
	"net/http"

	"gee"
)

func main() {
	fmt.Println("enter")
	r := gee.New()
	fmt.Println("enter")
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})
	err := r.RUN(":9999")
	if err != nil {
		panic(err)
	}
}
