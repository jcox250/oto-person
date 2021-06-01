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
	flag.StringVar(&url, "url", "https://localhost:8000/oto/", "The URL")
	flag.StringVar(&method, "method", "", "the method to use e.g. Add|Show")
	flag.StringVar(&payload, "payload", "", "the payload to send")
	flag.StringVar(&certFile, "certfile", "cert.pem", "the cert file for the server")
	flag.StringVar(&clientCertFile, "client-certfile", "clientcert.pem", "the cert file for the client")
	flag.StringVar(&clientKeyFile, "client-keyfile", "clientkey.pem", "the key for the client")
	flag.Parse()
}

func main() {
	logger := loglvl.NewLogger(os.Stderr, false)

	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		logger.Error("msg", "failed to load key pair", "err", err)
		os.Exit(1)
	}

	certPool, err := newCertPool(certFile)
	if err != nil {
		logger.Error("msg", "failed to create certpool from certFile", "certFile", certFile, "err", err)
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
			logger.Error("msg", "failed to unmarshal payload to ShowRequest", "err", err)
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

func newCertPool(f string) (*x509.CertPool, error) {
	cert, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		return nil, err
	}
	return certPool, nil
}
