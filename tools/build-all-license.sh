#!/bin/bash

if [[ $1 != "" ]] ; then
    cd ..
    exec > $2
fi

echo
echo "DigitalRebar Provision License"
echo

cat LICENSE

#echo "TODO: Get the downloaded Assets Licenses"
echo
go mod vendor &>/dev/null
cd vendor
find . | grep LICENSE | while read line ; do
    echo
    echo "GO Package License: $line"
    echo
    cat $line
done
cd ..
rm -rf vendor &>/dev/null
echo
