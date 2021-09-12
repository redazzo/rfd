#!/usr/bin/env bash

function doWindowsBuild() {
  echo "Building Windows ..."
  GOOS=windows
  GOARCH=amd64
  go build -o bin/rfd.exe ./script/go_modules/rfd
  cp ./bin/rfd.exe .
}

function doLinuxBuild() {
  echo "Building Linux ..."
  GOOS=linux
  GOARCH=amd64
  go build -o bin/rfd ./script/go_modules/rfd
  cp ./bin/rfd .
}

case "$1" in
  windows)
    doWindowsBuild
    ;;
  linux)
    doLinuxBuild
    ;;
  all)
    doLinuxBuild
    doWindowsBuild
    ;;
  *)
    echo $"Usage: $0 {windows|linux|all}"
    exit 1
esac





