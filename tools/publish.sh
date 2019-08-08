#!/usr/bin/env bash

set -e

. tools/version.sh
version="v${MajorV}.${MinorV}.${PatchV}${Extra}"

DOIT=0
if [[ $version =~ ^v ]] ; then
    DOIT=1
fi
if [[ $DOIT == 0 ]] ;then
    echo "Not a publishing branch."
    exit 0
fi

mkdir -p rebar-catalog/drpcli/amd64/linux
cp bin/linux/amd64/drpcli rebar-catalog/drpcli/amd64/linux
mkdir -p rebar-catalog/drpcli/amd64/darwin
cp bin/darwin/amd64/drpcli rebar-catalog/drpcli/amd64/darwin
mkdir -p rebar-catalog/drpcli/amd64/windows
cp bin/windows/amd64/drpcli rebar-catalog/drpcli/amd64/windows
mkdir -p rebar-catalog/drpcli/arm64/linux
cp bin/linux/arm64/drpcli rebar-catalog/drpcli/arm64/linux
mkdir -p rebar-catalog/drpcli/armv7/linux
cp bin/linux/armv7/drpcli rebar-catalog/drpcli/armv7/linux

