#!/usr/bin/env bash
. tools/version.sh >/dev/null
[[ $GOOS ]] || export GOOS=$(go env GOOS)
[[ $GOARCH ]] || export GOARCH=$(go env GOARCH)
[[ $GOARCH = arm ]] && export GOARM=7
binpath="$PWD/bin/$GOOS/${GOARCH}${GOARM+v${GOARM}}"
export PATH="$PWD/$binpath:$PATH"
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
          -X github.com/digitalrebar/provision/v4.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` \
          -X github.com/digitalrebar/provision/v4.GitHash=$GITHASH"
set -e
cd "$1"
if grep -qs 'go:generate' *; then
    go generate
fi
go build -ldflags "$VERFLAGS" -o "$binpath/$exename"
