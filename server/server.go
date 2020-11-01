package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
)

var
	REQUEST_COUNT = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_app_requests_count",
		Help: "Total App HTTP Requests count.",
	})

func main() {
	// Start the application
	startMyApp()
}

func startMyApp() {
	router := mux.NewRouter()
	router.HandleFunc("/birthday/{name}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		greetings := fmt.Sprintf("Happy Birthday %s :)", name)
		rw.Write([]byte(greetings))

		REQUEST_COUNT.Inc()
	}).Methods("GET")

	log.Println("Starting the application server...")
	router.Path("/metrics").Handler(promhttp.Handler())

	caCert, err := ioutil.ReadFile("./certs/ca.crt")
	if err != nil{
		fmt.Println(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      "0.0.0.0:9443",
		TLSConfig: tlsConfig,
		Handler: router,
	}
	fmt.Println("starting server")
	_ = server.ListenAndServeTLS("./certs/server.crt", "./certs/server.key")
}