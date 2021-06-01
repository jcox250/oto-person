package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jcox250/loglvl"
	"github.com/jcox250/oto-person/cache"
	"github.com/jcox250/oto-person/clients"
	servergen "github.com/jcox250/oto-person/gen/server"
	"github.com/jcox250/oto-person/service"
	"github.com/pacedotdev/oto/otohttp"
)

var (
	debug          bool
	port           int
	metricsPort    int
	certFile       string
	clientCertFile string
	keyFile        string
	redisAddr      string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "enables debug logs")
	flag.IntVar(&port, "port", 8000, "port the person service runs on")
	flag.IntVar(&metricsPort, "metrics-port", 9000, "port the metrics server runs on")
	flag.StringVar(&certFile, "certfile", "cert.pem", "certificate PEM file")
	flag.StringVar(&clientCertFile, "client-certfile", "clientcert.pem", "client certificate PEM file")
	flag.StringVar(&keyFile, "keyfile", "key.pem", "key PEM file")
	flag.StringVar(&redisAddr, "redis-addr", "localhost:6379", "address for redis")
	flag.Parse()
}

func main() {
	logger := loglvl.NewLogger(os.Stderr, debug)
	logger.Info("msg", "service config", "debug", debug, "port", port, "metrics-port", metricsPort)
	logger.Info("msg", "tls config", "keyfile", keyFile, "certfile", certFile, "client-certfile", clientCertFile)

	clientCertPool, err := loadCert(clientCertFile)
	if err != nil {
		logger.Error("msg", "failed to read client cert file", "err", err)
		os.Exit(1)
	}

	rc := clients.NewRedisClient(redisAddr, logger)
	personCache := cache.NewPersonCache(logger, rc)
	personService := service.New(logger, personCache)

	otoHandler := otohttp.NewServer()
	servergen.RegisterPersonService(otoHandler, personService)
	logger.Debug("msg", "registered person service")

	mux := http.NewServeMux()
	mux.Handle("/oto/", otoHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		TLSConfig: &tls.Config{
			ClientCAs:  clientCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	logger.Info("msg", fmt.Sprintf("serving on port %d...", port))
	log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}

func loadCert(f string) (*x509.CertPool, error) {
	clientCert, err := os.ReadFile(clientCertFile)
	if err != nil {
		return nil, err
	}
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCert)
	return clientCertPool, nil
}
