package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestServeProxy(t *testing.T) {
	// Start the proxy server
	proxyAddress := "localhost:8081"

	go serveProxy(&Config{Listen: proxyAddress, Rules: []Rule{
		{
			Domains: []string{
				".*",
			},
			Proxy: "local",
		},
	}})

	// Wait for server to start up
	for {
		conn, err := net.Dial("tcp", proxyAddress)
		if err == nil {
			conn.Close()
			break
		}
	}

	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", proxyAddress))
	if err != nil {
		log.Fatal(err)
	}

	// Create a transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	// Create a client
	client := &http.Client{
		Transport: transport,
	}

	for _, testUrl := range []string{"http://ifconfig.me", "https://ifconfig.me"} {
		log.Printf("[%s] sending GET", testUrl)
		resp, err := client.Get(testUrl)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		// Check response status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code %v", resp.StatusCode)
		}

		// Print out the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("[%s] Response body: %s\n", testUrl, string(body))
	}
}
