# Copyright 2021 The NMPolicy Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#


set -e

options=$(getopt --options "" \
    --long build,lint,unit-test,integration-test,help\
    -- "${@}")
eval set -- "$options"
while true; do
    case "$1" in
    --build)
        OPT_BUILD=1
        ;;
    --lint)
        OPT_LINT=1
        ;;
    --unit-test)
        OPT_UTEST=1
        ;;
    --integration-test)
        OPT_ITEST=1
        ;;
    --help)
        set +x
        echo "$0 [--build] [--lint] [--unit-test] [--integration-test]"
        exit
        ;;
    --)
        shift
        break
        ;;
    esac
    shift
done

if [ -z "${OPT_BUILD}" ] && [ -z "${OPT_LINT}" ] && [ -z "${OPT_UTEST}" ] && [ -z "${OPT_ITEST}" ]; then
    OPT_BUILD=1
    OPT_LINT=1
    OPT_UTEST=1
    OPT_ITEST=1
fi

if [ -n "${OPT_BUILD}" ]; then
    go build -o ./.out/nmpolicyctl ./cmd/nmpolicyctl/...
fi

if [ -n "${OPT_LINT}" ]; then
    golangci_lint_version=v1.42.1
    GOFLAGS=-mod=mod go run github.com/golangci/golangci-lint/cmd/golangci-lint@$golangci_lint_version run
fi

if [ -n "${OPT_UTEST}" ]; then
    go test -v ./nmpolicy/...
fi

if [ -n "${OPT_ITEST}" ]; then
    go test -v ./tests/...
fi
