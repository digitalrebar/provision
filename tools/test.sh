#!/usr/bin/env bash
PATH="$PWD/bin/$(go env GOOS)/$(go env GOARCH):$PATH"

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
    time go test -timeout 30m -race -covermode=atomic -coverpkg=$packages -coverprofile="profile${i}.txt" "$d" || FAILED=true
    i=$((i+1))
done
