#!/usr/bin/env bash

set -o errexit -o nounset -o pipefail

go install github.com/nektos/act@latest

act $@
