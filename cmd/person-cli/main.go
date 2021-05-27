package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jcox250/loglvl"
	"github.com/jcox250/oto-person/gen/client"
)

var (
	url            string
	method         string
	payload        string
	certFile       string
	clientCertFile string
	clientKeyFile  string
)

func init() {
	flag.StringVar(&url, "url", "", "the url for the person service")
	flag.StringVar(&method, "method", "", "the method to use e.g. Add|Show")
	flag.StringVar(&payload, "payload", "", "the payload to send")
	flag.StringVar(&certFile, "certfile", "", "the cert file for the server")
	flag.StringVar(&clientCertFile, "client-certfile", "", "the cert file for the client")
	flag.StringVar(&clientKeyFile, "keyfile", "", "the key for the client")
	flag.Parse()
}

func main() {
	logger := loglvl.NewLogger(os.Stderr, false)

	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		logger.Error("msg", "failed to load key pair", "err", err)
		os.Exit(1)
	}

	cert, err := os.ReadFile(certFile)
	if err != nil {
		logger.Error("msg", "failed to read certfile", "err", err)
		os.Exit(1)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		logger.Error("msg", "unable to parse cert")
		os.Exit(1)
	}

	c := client.New(url)
	c.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{clientCert},
		},
	}
	personClient := client.NewPersonService(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch method {
	case "Add":
		addRequest := client.AddRequest{}
		if err := json.Unmarshal([]byte(payload), &addRequest); err != nil {
			logger.Error("msg", "failed to unmarshal payload to AddRequest", "err", err)
			os.Exit(1)
		}

		resp, err := personClient.Add(ctx, addRequest)
		if err != nil {
			logger.Error("msg", "Add request failed", "err", err)
			os.Exit(1)
		}
		fmt.Println(resp)
	case "Show":
		showRequest := client.ShowRequest{}
		if err := json.Unmarshal([]byte(payload), &showRequest); err != nil {
			logger.Error("msg", "failed to unmarshal payload to AddRequest", "err", err)
			os.Exit(1)
		}

		resp, err := personClient.Show(ctx, showRequest)
		if err != nil {
			logger.Error("msg", "Add request failed", "err", err)
			os.Exit(1)
		}
		fmt.Println(resp)
	default:
		log.Fatalf("invalid method %q, valid methods are Add, Show", method)
	}

}
