package main

import (
	"flag"
	"os"

	"log"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "reload" {
			serviceReload()
		} else if os.Args[1] == "start" {
			serviceStart()
		} else if os.Args[1] == "stop" {
			serviceStop()
		}
	}
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	log.Printf("Configuration loaded: %v", config)

	err = serveProxy(config)
	if err != nil {
		panic(err)
	}
}
