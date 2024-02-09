#!/bin/bash

set -euo pipefail

cd $(dirname $0)/..
go get github.com/xlab/c-for-go
go mod download
GOBIN=$(pwd)/bin go install github.com/xlab/c-for-go
