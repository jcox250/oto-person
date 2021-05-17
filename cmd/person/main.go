package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jcox250/loglvl"
	servergen "github.com/jcox250/oto-person/gen/server"
	personservice "github.com/jcox250/oto-person/person_service"
	"github.com/pacedotdev/oto/otohttp"
)

var (
	debug       bool
	port        int
	metricsPort int
	certFile    string
	keyFile     string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "enables debug logs")
	flag.IntVar(&port, "port", 8000, "port the person service runs on")
	flag.IntVar(&metricsPort, "metrics-port", 9000, "port the metrics server runs on")
	flag.StringVar(&certFile, "certfile", "cert.pem", "certificate PEM file")
	flag.StringVar(&keyFile, "keyfile", "key.pem", "key PEM file")
	flag.Parse()
}

func main() {
	logger := loglvl.NewLogger(os.Stderr, debug)
	logger.Info("msg", "service config", "debug", debug, "port", port, "metrics-port", metricsPort)
	logger.Info("msg", "tls config", "keyfile", keyFile, "certfile", certFile)

	personService := personservice.New(logger)
	server := otohttp.NewServer()
	servergen.RegisterPersonService(server, personService)
	logger.Debug("msg", "registered person service")

	http.Handle("/oto/", server)
	logger.Info("msg", fmt.Sprintf("serving on port %d...", port))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certFile, keyFile, nil))
}
