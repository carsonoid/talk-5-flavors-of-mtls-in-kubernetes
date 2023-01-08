package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("USAGE: %s CERTFILE KEYFILE CAFILE\n", os.Args[0])
		os.Exit(1)
	}

	certFile := os.Args[1]
	keyFile := os.Args[2]
	caFile := os.Args[3]
	// END ARGS OMIT

	// START CA OMIT
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		panic(err)
	}
	caCertPool.AppendCertsFromPEM(caCert)
	// END CA OMIT

	// START TLS OMIT
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	// tlsConfig.BuildNameToCertificate()
	// END TLS OMIT

	// START SERVE OMIT
	server := &http.Server{
		Addr: ":8443",

		// use our custom TLS config with the ca, cert, and key
		// we *have* to use this to support a custom trusted CA pool
		TLSConfig: tlsConfig,

		// set all routes to simply print request info
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			names := r.TLS.VerifiedChains[0][0].DNSNames
			fmt.Printf("Verified Client Names: %v\n", names)
			fmt.Println("Request Headers:", r.Header)
			fmt.Fprintf(w, "Hello %s! Thanks for using mTLS! ",
				strings.Join(names, "/"))
		}),
	}

	fmt.Println("Starting HTTPS server on :8443")
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
	// END SERVE OMIT
}
