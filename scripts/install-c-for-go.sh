#!/bin/bash

set -euo pipefail

cd $(dirname $0)/..

GOBIN=$(pwd)/bin go install github.com/xlab/c-for-go@latest

# Comment out the above line and uncomment the line below to use the
# DigitalOcean fork of c-for-go. (This version can sometimes include
# bug fixes that are unreleased in the official version.)
# GOBIN=$(pwd)/bin go install github.com/digitalocean-labs/c-for-go@latest
