#!/usr/bin/env bash
tag_re='v([0-9]+)\.([0-9]+)\.([0-9]+)-(.+-)?([0-9]+)-(g[0-9a-f]+)'
TAG=$(git describe --tags --long --match 'v[0-9]*.[0-9]*.[0-9]*' --abbrev=1000)
if ! [[ $TAG =~ $tag_re ]]; then
    echo "Failed to find a semantic version tag!"
    echo "Add one with `git tag`"
    exit 1
fi >&2
MajorV="${BASH_REMATCH[1]}"
MinorV="${BASH_REMATCH[2]}"
PatchV="${BASH_REMATCH[3]}"
PRE="${BASH_REMATCH[4]%-}"
AHEAD="${BASH_REMATCH[5]}"
GITHASH="${BASH_REMATCH[6]}"
if [[ ! $PRE && $AHEAD != 0 ]]; then
   PRE="dev"
   PatchV=$((PatchV + 1))
fi

[[ $PRE ]] && Extra="-${PRE}.${AHEAD}"
echo "Version = v${MajorV}.${MinorV}.${PatchV}${Extra}+${GITHASH}"
