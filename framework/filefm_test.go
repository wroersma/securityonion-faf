// Copyright 2019 Jason Ertel (jertel). All rights reserved.
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
	"flag"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"os"
	"securityonion-faf/config"
	"testing"
	"time"
)

func TestInitLogging(tester *testing.T) {
	testFile := "/tmp/faf_test.log"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(testFile)
	file, err := InitLogging(testFile, "debug")
	if err != nil {
		tester.Errorf("expected no errors")
	}
	if file == nil {
		tester.Errorf("expected non-nil log file")
	}
}
func TestGetMD5Hash(tester *testing.T) {
	testFile := "/tmp/test.txt"
	err := ioutil.WriteFile(testFile, []byte("test"), 0600)
	testhash, err := GetMD5Hash(testFile)
	if err != nil {
		tester.Errorf("expected no errors")
	}
	testnohash, err := GetMD5Hash(testFile + "na")
	if testnohash != "" && err != nil {
		tester.Errorf("expected no errors")
	}
	if testhash == "d8e8fca2dc0f896fd7cb4cb0031ba249" {
		tester.Errorf("expected non-nil log file")
	}
}

func TestProcessFileDir(tester *testing.T) {
	configFilename := flag.String("c", "../faf.json", "Configuration file, in JSON format")
	testFile := zeekcom + "test.txt"
	err := ioutil.WriteFile(testFile, []byte("test"), 0600)
	if err != nil {
		tester.Errorf("expected no errors")
	}
	var fafs string
	fafs = "/tmp/faf.json"
	errrr := ioutil.WriteFile(fafs, []byte("{\"logLevel\":\"info\",\"logFilename\":\"faf.log\",\"redis\":{\"serverUrl\":\"\",\"serverPort\":\"\",\"serverPassword\":\"\"}}"), 0600)
	if errrr != nil {
		tester.Errorf("expected no errors")
	}
	configFilename1 := flag.String("e", fafs, "Configuration file, in JSON format")
	buildTime, err := time.Parse("2006-01-02T15:04:05", "unknown")
	cffg, err := config.LoadConfig(*configFilename1, "Test", buildTime)
	if err != nil {
		tester.Errorf("expected no errors")
	}
	if cffg != nil && cffg.Redis.ServerUrl != "" {
		tester.Errorf("expected no errors")
	}
	cfg, err := config.LoadConfig(*configFilename, "Test", buildTime)
	if err != nil {
		tester.Errorf("expected no errors")
	}
	ProcessFileDir(cfg)
	testFile2 := strelkaunproc + "test.txt"
	testhash, err := GetMD5Hash(testFile2)
	if err != nil {
		tester.Errorf("expected no errors")
	}
	if testhash != "098f6bcd4621d373cade4e832627b4f6" {
		tester.Errorf("expected non-nil log file")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.ServerUrl + ":" + cfg.Redis.ServerPort,
		Password: cfg.Redis.ServerPassword,
		DB:       0,
	})
	val, err := rdb.Get(ctx, testhash).Result()
	if err != nil {
		tester.Errorf("Expected to find hash value in redis.")
	}
	if val != "md5hash" {
		tester.Errorf("Expected to find md5hash value from redis.")
	}
	val2, err := rdb.Get(ctx, "d8e8fca2dc0f896fd7cb4cb0031ba248").Result()
	if err == nil {
		tester.Errorf("Expected to error on this value search with redis.")
	}
	if val2 != "" {
		tester.Errorf("Expected to find md5hash value from redis.")
	}
}
