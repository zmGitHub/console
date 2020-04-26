#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/custm-chat main.go