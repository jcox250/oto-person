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
)

func init() {
	flag.BoolVar(&debug, "debug", false, "enables debug logs")
	flag.IntVar(&port, "port", 8000, "port the person service runs on")
	flag.IntVar(&metricsPort, "metrics-port", 9000, "port the metrics server runs on")
	flag.Parse()
}

func main() {
	logger := loglvl.NewLogger(os.Stderr, debug)
	logger.Info("msg", "service config", "debug", debug, "port", port, "metrics-port", metricsPort)

	personService := personservice.New(logger)
	server := otohttp.NewServer()
	servergen.RegisterPersonService(server, personService)
	logger.Debug("msg", "registered person service")

	http.Handle("/oto/", server)
	logger.Info("msg", fmt.Sprintf("serving on port %d...", port))
	//log.Fatal(http.ListenAndServe(":8080", nil))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
