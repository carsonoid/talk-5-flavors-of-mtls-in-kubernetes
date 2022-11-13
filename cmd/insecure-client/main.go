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

	for range time.Tick(time.Second) {
		resp, err := client.Get(target)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(resp.Header)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("< GOT:", string(body))
	}
}
