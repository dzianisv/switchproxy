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
	go serveProxy(&Config{Listen: "localhost:8080", Rules: []Rule{
		{
			Domains: []string{
				".*",
			},
			Proxy: "local",
		},
	}})

	// Wait for server to start up
	for {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err == nil {
			conn.Close()
			break
		}
	}

	proxyURL, err := url.Parse("http://localhost:8080")
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

	resp, err := client.Get("https://ifconfig.me")
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
	fmt.Printf("Response body: %s\n", string(body))
}
