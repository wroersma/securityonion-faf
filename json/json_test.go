// Copyright 2019 Jason Ertel (jertel). All rights reserved.
// Copyright 2020-2021 Security Onion Solutions, LLC. All rights reserved.
//
// This program is distributed under the terms of version 2 of the
// GNU General Public License.  See LICENSE for further details.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

package json

import (
	"os"
	"testing"
)

func TestJson(tester *testing.T) {
	testFile := "/tmp/faf_test.go.json"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			tester.Errorf("Unexpected error while removing the test file.")
		}
	}(testFile)
	obj := make(map[string]string)
	obj["MyKey"] = "MyValue"
	err := WriteJsonFile(testFile, obj)
	if err != nil {
		tester.Errorf("Unexpected write error")
	}
	errr := WriteJsonFile("test"+testFile, obj)
	if errr == nil {
		tester.Errorf("Expected write error")
	}
	obj = make(map[string]string)
	err = LoadJsonFile(testFile, &obj)
	if err != nil {
		tester.Errorf("unexpected load error")
	}
	err = LoadJson([]byte("{\"test: \"test\"}"), &obj)
	if err == nil {
		tester.Errorf("Expected a load error")
	}
	err = LoadJson([]byte("{\"test\": false}"), &obj)
	if err == nil {
		tester.Errorf("Expected load error")
	}
	if obj["MyKey"] != "MyValue" {
		tester.Errorf("expected value %s but got %s", "MyValue", obj["MyKey"])
	}
}
