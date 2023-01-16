package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("USAGE: %s CERTFILE KEYFILE CAFILE TARGET\n", os.Args[0])
		os.Exit(1)
	}

	certFile := os.Args[1]
	keyFile := os.Args[2]
	caFile := os.Args[3]
	target := os.Args[4]
	// END CLIENT ARGS OMIT

	// START CLIENT TLS OMIT
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		panic(err)
	}
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	// END CLIENT TLS OMIT

	// START CLIENT OMIT
	client := &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	// END CLIENT OMIT

	// START CLIENT LOOP OMIT
	for range time.Tick(time.Second) {
		resp, err := client.Get(target)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		names := resp.TLS.VerifiedChains[0][0].DNSNames
		fmt.Println("< Server certificate Names:", names)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("< Response Headers:", resp.Header)
		fmt.Println("< Response Body:", string(body))
	}
	// END CLIENT LOOP OMIT
}
