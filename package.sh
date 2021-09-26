#!/usr/bin/env bash

if [ -d "./package" ]
then
  rm -rvf ./package
  mkdir package
fi

if [ -e "./rfd.tar" ]
then
  rm -vf rfd.tar
fi

/bin/bash ./build.sh

cp -v ./rfd ./package
cp -rvf ./template ./package
cp -v ./config.yml ./package
tar -cvf rfd.tar ./package
