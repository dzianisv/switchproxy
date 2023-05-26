package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"
)

func TestServeProxy(t *testing.T) {
	// Start the proxy server
	go serveProxy(&Config{Listen: "localhost:8080", Rules: []Rule{
		{
			Domains: []string{
				"*",
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

	// Configure the transport to use the proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
		}),
		DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			return net.Dial("tcp", "ifconfig.me:80")
		},
	}

	// Make the request using the transport
	client := &http.Client{Transport: transport}
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
