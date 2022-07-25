#!/bin/bash
echo "Building Windows amd64..."
GOOS=windows GOARCH=amd64 go build -o Build/ytac-amd64.exe ytac.go
echo "Building Windows arm64..."
GOOS=windows GOARCH=arm64 go build -o Build/ytac-arm64.exe ytac.go
echo "Building Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o Build/ytac-linux-amd64 ytac.go
echo "Building Linux arm64..."
GOOS=linux GOARCH=arm64 go build -o Build/ytac-linux-arm64 ytac.go
echo "Building MacOS amd64..."
GOOS=darwin GOARCH=amd64 go build -o Build/ytac-macos-amd64 ytac.go
echo "Building MacOS arm64..."
GOOS=darwin GOARCH=arm64 go build -o Build/ytac-macos-arm64 ytac.go
echo "Done."
