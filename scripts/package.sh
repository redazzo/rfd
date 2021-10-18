#!/usr/bin/env bash

if [ -d "../build" ]
then
  rm -rvf ../build
fi

mkdir ../build

if [ -e "../rfd.tar" ]
then
  rm -vf rfd.tar
fi

/bin/bash ./build.sh

cp -v ../rfd ../build
cp -rvf ../template ../build
cp -v ../config.yml ../build
tar -cvf rfd.tar ../build
