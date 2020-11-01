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
	cert, err := ioutil.ReadFile("./certs/ca.crt")
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair("./certs/client.crt", "./certs/client.key")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}

	client := http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
				Certificates:[]tls.Certificate{certificate},
			},
		},
	}

	resp, err := client.Get("https://exampletest.com:9443/metrics")
	if err != nil{
		log.Fatal(err)
	}
	byteArr, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(string(byteArr))
	fmt.Println(resp.StatusCode)

}


