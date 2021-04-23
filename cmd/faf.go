// Copyright 2021 Wyatt Roersma (wroersma). All rights reserved.
// Copyright 2020-2021 Security Onion Solutions, LLC. All rights reserved.
//
// This program is distributed under the terms of version 2 of the
// GNU General Public License.  See LICENSE for further details.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
package main

import (
	"flag"
	"fmt"
	"github.com/apex/log"
	"os"
	"securityonion-faf/config"
	"securityonion-faf/framework"
	"time"
)

var (
	BuildVersion = "unknown"
	BuildTime    = "unknown"
)

func main() {
	// Load config file.
	configFilename := flag.String("c", "faf.json", "Configuration file, in JSON format")
	flag.Parse()
	buildTime, err := time.Parse("2006-01-02T15:04:05", BuildTime)
	if err != nil {
		fmt.Printf("Unable to parse build time; reason=%s\n", err.Error())
	}
	cfg, err := config.LoadConfig(*configFilename, BuildVersion, buildTime)
	if err == nil {
		//Start logging to a file if possible.
		logFile, _ := framework.InitLogging(cfg.LogFilename, cfg.LogLevel)
		defer func(logFile *os.File) {
			err := logFile.Close()
			if err != nil {
				log.Fatal("Log file not closed properly!")
			}
		}(logFile)
		log.WithFields(log.Fields{
			"version":   cfg.Version,
			"buildTime": cfg.BuildTime,
		}).Info("Version Information")
	}
	if cfg != nil {
		if err == nil && cfg.Redis != nil {
			// Start the processing of new files extracted by zeek.
			framework.ProcessFileDir(cfg)
		} else if cfg.Redis == nil {
			log.Fatal("Redis information not set in faf.json")
		}
	}
}
