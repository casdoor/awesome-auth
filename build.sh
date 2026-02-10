#!/bin/bash
export GOPROXY="https://proxy.golang.org,direct"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server .
