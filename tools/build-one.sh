#!/usr/bin/env bash
. tools/version.sh >/dev/null
[[ $BUILDSTAMP ]] || export BUILDSTAMP=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
[[ $GOOS ]] || export GOOS=$(go env GOOS)
[[ $GOARCH ]] || export GOARCH=$(go env GOARCH)
[[ $GOARCH = arm ]] && export GOARM=7
binpath="$PWD/bin/$GOOS/${GOARCH}${GOARM+v${GOARM}}"
export PATH="$PWD/bin/$(go version | awk '{ print $4 }'):$PATH"
exename="${1##*/}"
[[ $GOOS = "windows" ]] && exename="${exename}.exe"
mkdir -p "$binpath"
export GO111MODULE=on
export CGO_ENABLED=0
export VERFLAGS="-s -w \
          -X github.com/digitalrebar/provision/v4.RSMajorVersion=$MajorV \
          -X github.com/digitalrebar/provision/v4.RSMinorVersion=$MinorV \
          -X github.com/digitalrebar/provision/v4.RSPatchVersion=$PatchV \
          -X github.com/digitalrebar/provision/v4.RSExtra=$Extra \
          -X github.com/digitalrebar/provision/v4.BuildStamp=$BUILDSTAMP"
set -e
cd "$1"
if [[ $TRAVIS = true ]]; then
    # Sigh.  Work around some rate-limiting hoodoo, hopefully
    for i in 1 2 3 4 5; do
        go mod download && break
        sleep $i
    done
fi
if grep -qs 'go:generate' *; then
    go generate
fi
go build -ldflags "$VERFLAGS" -o "$binpath/$exename"
