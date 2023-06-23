#!/usr/bin/env bash

function doWindowsBuild() {
  echo "Building Windows binary. Binaries are in the /bin directory  ..."
  GOOS=windows
  GOARCH=amd64
  go build -o ./bin/rfd.exe ./cmd/rfd
  echo "Done"
}

function doLinuxBuild() {
  echo "Building Linux binary. Binaries are in the /bin directory ..."
  GOOS=linux
  GOARCH=amd64
  go build -o ./bin/rfd ./cmd/rfd
  echo "Done"
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





