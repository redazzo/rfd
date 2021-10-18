#!/usr/bin/env bash

function doWindowsBuild() {
  echo "Building Windows binary ..."
  GOOS=windows
  GOARCH=amd64
  go build -o ../bin/rfd.exe ../cmd/rfd
  cp ../bin/rfd.exe ..
}

function doLinuxBuild() {
  echo "Building Linux binary ..."
  GOOS=linux
  GOARCH=amd64
  go build -o ../bin/rfd ../cmd/rfd
  cp ../bin/rfd ..
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
    echo $"Defaulting to linux ..."
    doLinuxBuild
    exit 1
esac





