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
	"github.com/apex/log"
)

type RedisConfig struct {
	ServerUrl      string `json:"serverUrl"`
	ServerPort     string `json:"serverPort"`
	ServerPassword string `json:"serverPassword"`
}

func (config *RedisConfig) Verify() error {
	var err error
	if err == nil && config.ServerUrl == "" {
		log.Fatal("Redis.ServerUrl configuration value is required")
	}
	if err == nil && config.ServerPort == "" {
		log.Fatal("Redis.ServerPort configuration value is required")
	}
	return err
}
