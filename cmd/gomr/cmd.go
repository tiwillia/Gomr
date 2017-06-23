package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/tiwillia/gomr/pkg/gomr"
)

// This is where we start heh
func main() {
	configPath := flag.String("config", "./gomr.yaml", "Path to config file")
	flag.Parse()
	glog.Infoln("Starting irc bot...")

	// Read configuration
	config, dbConfig := gomr.GetConfiguration(*configPath)

	gomrService, err := gomr.NewGomrService(&config, &dbConfig)
	if err != nil {
		glog.Fatalf("Unable to create Gomr service: %s", err)
	}

	err = gomrService.Run()
	if err != nil {
		glog.Fatalf("Error encountered running Gomr service: %s", err)
	}
}
