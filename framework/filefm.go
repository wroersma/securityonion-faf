// Package framework
// Copyright 2021 Wyatt Roersma (wroersma). All rights reserved.
// Copyright 2020-2021 Security Onion Solutions, LLC. All rights reserved.
//
// This program is distributed under the terms of version 2 of the
// GNU General Public License.  See LICENSE for further details.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
package framework

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
	"github.com/go-redis/redis/v8"
	"io"
	"io/ioutil"
	"os"
	"securityonion-faf/config"
)

var (
	zeekcom = "/nsm/zeek/extracted/complete/"
	ctx     = context.Background()
)

func InitLogging(logFilename string, logLevel string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		log.SetHandler(logfmt.New(logFile))
	} else {
		log.WithError(err).WithField("logfile", logFilename).Error("Failed to create log file, " +
			"using console instead")
		log.SetHandler(text.New(os.Stdout))
	}
	log.SetLevelFromString(logLevel)
	return logFile, err
}

func GetMD5Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ProcessFileDir(cfg *config.Config) {
	// Use faf config to set redis options
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.ServerUrl + ":" + cfg.Redis.ServerPort,
		Password: cfg.Redis.ServerPassword,
		DB:       0,
	})
	files, err := ioutil.ReadDir("/nsm/zeek/extracted/complete")
	if err != nil {
	}
	for _, f := range files {
		var newfile string
		newfile, err = GetMD5Hash(zeekcom + f.Name())
		if err != nil {
			log.Info("Error hashing the extracted zeek file.")
		}
		val, err := rdb.Get(ctx, newfile).Result()
		if err != nil {
			oldLocation := zeekcom + f.Name()
			newLocation := "/nsm/strelka/unprocessed/" + f.Name()
			err := os.Rename(oldLocation, newLocation)
			if err != nil {
			}
			err = rdb.Set(ctx, newfile, "md5hash", 0).Err()
			// Print out error if detected
			if err != nil {
				log.Info("Failed getting message from redis")
			}
			log.Info("Moved file " + f.Name() + " to be processed with md5hash: " + newfile)
		} else {
			// Remove this file
			os.Remove(zeekcom + f.Name())
			log.Info("Removing duplicate file: " + f.Name() + " with md5hash: " + newfile)
		}
		if val == "" {
		}
	}
}
