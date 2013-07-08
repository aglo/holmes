#!/bin/bash

GOPATH=`pwd`

git submodule init
git submodule update

go install holmes
mkdir data
cp conf/holmes.conf bin
