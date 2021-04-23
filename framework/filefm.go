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
	zeekcom       = "/nsm/zeek/extracted/complete/"
	ctx           = context.Background()
	strelkaunproc = "/nsm/strelka/unprocessed/"
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
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
	files, err := ioutil.ReadDir(zeekcom[:len(zeekcom)-1])
	if err != nil {
		log.Error("Failed reading directory " + zeekcom)
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
			newLocation := strelkaunproc + f.Name()
			err := os.Rename(oldLocation, newLocation)
			if err != nil {
			}
			err = rdb.Set(ctx, newfile, "md5hash", 0).Err()
			// Print out error if detected
			if err != nil {
				log.Info("Failed getting message from redis.")
			}
			log.Info("Moved file to be processed with md5hash: " + newfile + " with file name: " + f.Name())
		} else {
			// Remove this file
			err := os.Remove(zeekcom + f.Name())
			if err != nil {
				log.Error("Failed removing duplicate  with md5hash: " + newfile + " with file name: " + f.Name())
			}
			log.Info("Removing duplicate  with md5hash: " + newfile + " with file name: " + f.Name())
		}
		if val == "" {
		}
	}
}
