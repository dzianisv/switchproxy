package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"log"

	"github.com/elazarl/goproxy"
)

func getTLSConfig(host string, ctx *goproxy.ProxyCtx) (*tls.Config, error) {
	return &tls.Config{InsecureSkipVerify: true}, nil
}
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

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.Tr.Dial = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial

	// For other requests, use the appropriate upstream proxy
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		rule := getMatchingRule(config, r.URL.Host)
		log.Printf("DoFunc: %v %v", rule, r.URL.Host)
		if rule != nil && rule.Proxy != "local" {
			upstreamProxy, err := url.Parse(rule.Proxy)
			if err != nil {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusInternalServerError,
					fmt.Sprintf("Error parsing proxy URL: %v", err))
			}
			log.Printf("Using upstream proxy \"%s\" for \"%s\"", upstreamProxy, r.URL.Host)
			proxy.Tr.Proxy = http.ProxyURL(upstreamProxy)
		} else {
			proxy.Tr.Proxy = nil
		}
		return r, nil
	})

	log.Printf("Starting proxy server on %s", ":8080")
	if err := http.ListenAndServe(":8080", proxy); err != nil {
		log.Fatalf("Error starting proxy server: %v", err)
	}
}
