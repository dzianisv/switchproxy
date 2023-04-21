package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"

	"log"

	"github.com/elazarl/goproxy"
)

func getMatchingRule(config Config, host string) *Rule {
	for _, rule := range config.Rules {
		for _, domain := range rule.Domains {

			matched, err := regexp.MatchString(domain, host)
			if err != nil {
				return nil
			}

			if matched {
				return &rule
			}
		}
	}
	return nil
}

func ConnectProxy(network string, addr string, proxyURL string) (net.Conn, error) {
	upstreamProxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	proxyConn, err := net.Dial("tcp", upstreamProxy.Host)
	if err != nil {
		return nil, err
	}

	connectReq := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: addr},
		Host:   addr,
		Header: make(http.Header),
	}

	if upstreamProxy.User != nil {
		password, _ := upstreamProxy.User.Password()
		auth := base64.StdEncoding.EncodeToString([]byte(upstreamProxy.User.Username() + ":" + password))
		connectReq.Header.Set("Proxy-Authorization", "Basic "+auth)
	}

	connectReq.Write(proxyConn)
	br := bufio.NewReader(proxyConn)
	resp, err := http.ReadResponse(br, connectReq)
	if err != nil {
		proxyConn.Close()
		return nil, err
	}
	if resp.StatusCode != 200 {
		proxyConn.Close()
		return nil, fmt.Errorf("non-200 status code from upstream proxy: %d", resp.StatusCode)
	}

	return proxyConn, nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	proxy := goproxy.NewProxyHttpServer()

	proxy.ConnectDial = func(network string, addr string) (net.Conn, error) {
		rule := getMatchingRule(config, addr)
		if rule != nil || rule.Proxy != "local" {
			return ConnectProxy(network, addr, rule.Proxy)
		} else {
			return net.Dial(network, addr)
		}
	}

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if req.Method == http.MethodConnect {
			return req, nil
		}

		rule := getMatchingRule(config, req.URL.Host)
		if rule != nil && rule.Proxy != "local" {
			upstreamProxy, err := url.Parse(rule.Proxy)
			if err != nil {
				log.Printf("Error parsing proxy URL for %s: %v", req.URL.Host, err)
				return req, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusInternalServerError, "Error parsing proxy URL")
			}
			proxy.Tr = &http.Transport{
				Proxy: http.ProxyURL(upstreamProxy),
			}
		} else {
			proxy.Tr, _ = http.DefaultTransport.(*http.Transport)
		}

		return req, nil
	})

	log.Printf("Starting proxy server on %s", config.Listen)
	log.Fatal(http.ListenAndServe(config.Listen, proxy))
}
