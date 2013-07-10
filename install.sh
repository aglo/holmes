#!/bin/bash

export GOPATH=`pwd`

git submodule init
git submodule update

./format.sh
go test holmes
go install holmes

if [ -e "data" ]
then
    :
else
    mkdir data
fi
cp conf/holmes.conf bin
