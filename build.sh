#!/bin/bash

env GOOS=windows GOARCH=amd64 go build -o binaries/integrity-win-amd64.exe integrity.go
env GOOS=darwin  GOARCH=amd64 go build -o binaries/integrity-mac-amd64 integrity.go
env GOOS=linux   GOARCH=amd64 go build -o binaries/integrity-linux-amd64 integrity.go
