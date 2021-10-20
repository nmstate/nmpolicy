#!/usr/bin/env bash

set -e

options=$(getopt --options "" \
    --long build,fmt,unit-test,integration-test,help\
    -- "${@}")
eval set -- "$options"
while true; do
    case "$1" in
    --build)
        OPT_BUILD=1
        ;;
    --fmt)
        OPT_FMT=1
        ;;
    --unit-test)
        OPT_UTEST=1
        ;;
    --integration-test)
        OPT_ITEST=1
        ;;
    --help)
        set +x
        echo "$0 [--build] [--fmt] [--unit-test] [--integration-test]"
        exit
        ;;
    --)
        shift
        break
        ;;
    esac
    shift
done

if [ -z "${OPT_BUILD}" ] && [ -z "${OPT_FMT}" ] && [ -z "${OPT_UTEST}" ] && [ -z "${OPT_ITEST}" ]; then
    OPT_BUILD=1
    OPT_FMT=1
    OPT_UTEST=1
    OPT_ITEST=1
fi

if [ -n "${OPT_BUILD}" ]; then
    go build -o ./.out/nmpolicy ./cmd/nmpolicy
fi

if [ -n "${OPT_FMT}" ]; then
        unformatted=$(gofmt -l ./nmpolicy ./tests)
        test -z "$unformatted" || (echo "Unformatted: $unformatted" && false)

        go get golang.org/x/tools/cmd/goimports
        unformatted=$(go run golang.org/x/tools/cmd/goimports -l --local "github.com/nmstate/nmpolicy" ./nmpolicy ./tests)
        go mod tidy

        test -z "$unformatted" || (echo "Unformatted imports: $unformatted" && false)
fi

if [ -n "${OPT_UTEST}" ]; then
    go test -v ./nmpolicy/...
fi

if [ -n "${OPT_ITEST}" ]; then
    go test -v ./tests/...
fi
