#!/usr/bin/env bash
GOOS=windows GOARCH=amd64 go build -o bin/rfd.exe ./script/go_modules/rfd
GOOS=linux GOARCH=amd64 go build -o bin/rfd ./script/go_modules/rfd
cp ./bin/rfd .