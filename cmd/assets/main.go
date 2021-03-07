package main

import (
	"fmt"
	"html"
	"net/http"
)

func main() {
	fmt.Println("hello from assets main")

	c := make(chan struct{})

	go func() {
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
			c <- struct{}{}
		})

		http.ListenAndServe(":4949", nil)
	}()

	<-c
}
