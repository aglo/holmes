#!/bin/bash

export GOPATH=`pwd`

if [ -e "data" ]
then
    :
else
    mkdir -p data/in_log
    mkdir -p data/out_log
fi

git submodule init
git submodule update

./format.sh
go test holmes
go install holmes

cp conf/holmes.conf bin