#!/bin/bash

# Copyright 2021 Wyatt Roersma (wroersma). All rights reserved.
#
# This program is distributed under the terms of version 2 of the
# GNU General Public License.  See LICENSE for further details.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

version=${1:-dev}
now=`date -u +%Y-%m-%dT%H:%M:%S`

go get ./...
go build -a -ldflags "-X main.BuildVersion=$version -X main.BuildTime=$now -extldflags '-static'" -tags faf -installsuffix faf cmd/faf.go