// Copyright 2021 Wyatt Roersma (wroersma). All rights reserved.
// Copyright 2020-2021 Security Onion Solutions, LLC. All rights reserved.
//
// This program is distributed under the terms of version 2 of the
// GNU General Public License.  See LICENSE for further details.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

package config

import (
	"securityonion-faf/json"
	"time"
)

type Config struct {
	Filename              string
	Version               string
	BuildTime             time.Time
	LoadTime              time.Time
	LogLevel              string       `json:"logLevel"`
	LogFilename           string       `json:"logFilename"`
	ShutdownGracePeriodMs int          `json:"shutdownGracePeriodMs"`
	Redis                 *RedisConfig `json:"redis"`
}

func LoadConfig(filename string, version string, buildTime time.Time) (*Config, error) {
	cfg := &Config{
		Version:               version,
		BuildTime:             buildTime,
		Filename:              filename,
		LoadTime:              time.Now(),
		LogLevel:              "info",
		LogFilename:           filename + ".log",
		ShutdownGracePeriodMs: 10000,
	}
	err := json.LoadJsonFile(cfg.Filename, cfg)
	if err == nil {
		if cfg.Redis != nil {
			err = cfg.Redis.Verify()
		}
	}
	return cfg, err
}
