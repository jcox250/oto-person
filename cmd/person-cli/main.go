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
	url      string
	method   string
	payload  string
	certFile string
)

func init() {
	flag.StringVar(&url, "url", "", "the url for the person service")
	flag.StringVar(&method, "method", "", "the method to use e.g. Add|Show")
	flag.StringVar(&payload, "payload", "", "the payload to send")
	flag.StringVar(&certFile, "certFile", "", "the cert file used by the server")
	flag.Parse()
}

func main() {
	logger := loglvl.NewLogger(os.Stderr, false)

	cert, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalf("failed to read certfile %s: %s", certFile, err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatal("unable to create certpoll from certFile")
	}

	c := client.New(url)
	c.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
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
