#!/usr/bin/env bash

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
    go build -o ./.out/nmpolicy ./cmd/nmpolicy
fi

if [ -n "${OPT_LINT}" ]; then
    golangci_lint_version=v1.42.1
    if [ ! -f $(go env GOPATH)/bin/golangci-lint ]; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $golangci_lint_version
    fi
    golangci-lint run
fi

if [ -n "${OPT_UTEST}" ]; then
    go test -v ./nmpolicy/...
fi

if [ -n "${OPT_ITEST}" ]; then
    go test -v ./tests/...
fi
