package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("USAGE: %s CERTFILE KEYFILE CAFILE\n", os.Args[0])
		os.Exit(1)
	}

	certFile := os.Args[1]
	keyFile := os.Args[2]
	caFile := os.Args[3]

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			commonName := r.TLS.VerifiedChains[0][0].Subject.CommonName
			fmt.Fprintf(w, "Hello %q from tls!", commonName)
		}),
	}

	fmt.Println("Starting HTTPS server on :8443")
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}
