#!/bin/bash

GOPATH=`pwd`

git submodule init
git submodule update

go install holmes
cp conf/holmes.conf bin