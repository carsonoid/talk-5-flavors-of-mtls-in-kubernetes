package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	target := os.Args[1]

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	fmt.Println("Starting HTTP client request loop")
	for range time.Tick(time.Second) {
		resp, err := client.Get(target)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("< Response Headers:", resp.Header)
		fmt.Println("< Response Body:", string(body))
	}
}
