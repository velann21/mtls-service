package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	io.WriteString(w, "Hello, world!\n")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hello",helloHandler)



	// load CA certificate file and add it to list of client CAs
	caCertFile, err := ioutil.ReadFile("/Users/singaravelannandakumar/go/src/awesomeProject1/certs/ca.crt")
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertFile)


	// serve on port 9090 of local host
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  certPool,
			MinVersion: tls.VersionTLS12,
		},
	}

	log.Println("Starting the application server...")
	// Listen to port 8080 and wait
	log.Fatal(server.ListenAndServeTLS( "/Users/singaravelannandakumar/go/src/awesomeProject1/certs/serverkey/server.crt", "/Users/singaravelannandakumar/go/src/awesomeProject1/certs/serverkey/server.key"))
}