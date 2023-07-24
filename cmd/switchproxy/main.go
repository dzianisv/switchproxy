package main

import (
	"flag"
	"fmt"

	"log"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	actionFlag := flag.String("service", "", "Service action: reload, start, stop, install")

	flag.Parse()
	var err error

	if len(*actionFlag) != 0 {
		if *actionFlag == "reload" {
			err = serviceReload()
		} else if *actionFlag == "start" {
			err = serviceStart()
		} else if *actionFlag == "stop" {
			err = serviceStop()
		} else if *actionFlag == "install" {
			err = serviceInstall()
		} else {
			err = fmt.Errorf("Unsupported action %s", *actionFlag)
		}

		if err != nil {
			log.Fatalf("failed to %s: %s", *actionFlag, err)
		}

		return
	}

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
