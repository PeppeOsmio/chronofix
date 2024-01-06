#!/bin/bash
cd src
GOOS=linux GOARCH=amd64 go build -o ../mediametaparser-linux-amd64 main.go
GOOS=linux GOARCH=386 go build -o ../mediametaparser-linux-386 main.go
GOOS=windows GOARCH=amd64 go build -o ../mediametaparser-windows-amd64.exe main.go
GOOS=windows GOARCH=386 go build -o ../mediametaparser-windows-386.exe main.go
GOOS=darwin GOARCH=amd64 go build -o ../mediametaparser-darwin-amd64 main.go
GOOS=darwin GOARCH=386 go build -o ../mediametaparser-darwin-386 main.go