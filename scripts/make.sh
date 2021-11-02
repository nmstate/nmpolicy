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
    --long build,headers:,lint,unit-test,integration-test,help\
    -- "${@}")
eval set -- "$options"
while true; do
    case "$1" in
    --build)
        OPT_BUILD=1
        ;;
    --headers)
        OPT_HEADERS=$2
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
        echo "$0 [--build] [--headers=[fix|check]] [--lint] [--unit-test] [--integration-test]"
        exit
        ;;
    --)
        shift
        break
        ;;
    esac
    shift
done

if [ -z "${OPT_BUILD}" ] && [ -z "${OPT_LINT}" ] && [ -z "${OPT_UTEST}" ] && [ -z "${OPT_ITEST}" ] && [ -z "${OPT_HEADERS}" ]; then
    OPT_BUILD=1
    OPT_LINT=1
    OPT_UTEST=1
    OPT_ITEST=1
    OPT_HEADERS=check
fi

if [ -n "${OPT_BUILD}" ]; then
    go build -o ./.out/nmpolicy ./cmd/nmpolicy
fi

if [ -n "${OPT_LINT}" ]; then
    golangci_lint_version=v1.42.1
    if [ ! -f $(go env GOPATH)/bin/golangci-lint ]; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $golangci_lint_version
    fi
    golangci-lint run
fi

if [ -n "${OPT_HEADERS}" ]; then
    if [ ! -f $(go env GOPATH)/bin/addlicense ]; then
        (
           pushd ~
           go install github.com/google/addlicense@latest
           popd
        )
    fi
    if [ "${OPT_HEADERS}" == "check" ]; then
        args=-check
    fi
    addlicense $args  -c "The NMPolicy Authors." .
fi

if [ -n "${OPT_UTEST}" ]; then
    go test -v ./nmpolicy/...
fi

if [ -n "${OPT_ITEST}" ]; then
    go test -v ./tests/...
fi
