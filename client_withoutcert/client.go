package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	cert, err := ioutil.ReadFile("/Users/singaravelannandakumar/go/src/awesomeProject1/certs/ca.crt")
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
	}

	certificate, err := tls.LoadX509KeyPair("/Users/singaravelannandakumar/go/src/awesomeProject1/certs/client/client.crt", "/Users/singaravelannandakumar/go/src/awesomeProject1/certs/client/client.key")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)
	client := http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
				Certificates:[]tls.Certificate{certificate},
			},
		},
	}

	// Request /hello over port 8080 via the GET method
	r, err := client.Get("https://velan.com:8080/hello")
	if err != nil {
		log.Fatal(err)
	}

	// Read the response body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Printf("%s\n", body)
}