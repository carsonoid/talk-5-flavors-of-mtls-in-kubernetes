package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("> Request Headers:", r.Header)
			fmt.Fprintln(w, "Hello!")
		}),
	}

	fmt.Println("Starting HTTP server on :8080")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
