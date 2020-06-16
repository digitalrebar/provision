#!/usr/bin/env bash
set -e
export PATH="$PWD/bin/$(go env GOOS)/$(go env GOARCH):$PATH"
export GO111MODULE=on

tools/build-one.sh cmds/drpcli
tools/build-one.sh cmds/incrementer

. tools/version.sh

if ! which dr-provision &>/dev/null; then
  drpcli catalog item download drp --version=v4.3.2
  tar -zxvf drp.zip "bin/$(go env GOOS)/$(go env GOARCH)/dr-provision"
  rm drp.zip
fi
if ! which dr-provision &>/dev/null; then
    echo "No dr-provision binary to run tests against"
    exit 1
fi
ver_re='v4\.[0-9]+\.[0-9]+.*'
if ! [[ $(dr-provision --version 2>&1) =~ $ver_re ]]; then
    echo "Make sure a dr-provision binary of at least v4.0.0 or later is in your PATH"
    exit 1
fi


echo Running with $(which dr-provision) version $BASH_REMATCH

packages="github.com/digitalrebar/provision/v4,\
github.com/digitalrebar/provision/v4/models,\
github.com/digitalrebar/provision/v4/plugin,\
github.com/digitalrebar/provision/v4/cli,\
github.com/digitalrebar/provision/v4/api,\
github.com/digitalrebar/provision/v4/agent\
"

i=0
for d in $(go list ./... 2>/dev/null | egrep -v 'cmds|test') ; do
    echo "----------- TESTING $d -----------"
    time go test -timeout 30m -race -covermode=atomic -coverpkg=$packages -coverprofile="profile${i}-c.txt" "$d" || FAILED=true
    i=$((i+1))
done
rm "bin/$(go env GOOS)/$(go env GOARCH)/dr-provision" || :
[[ ! $FAILED ]]
