#!/usr/bin/env bash

# Cross-Compile golang
env GOOS=linux GOARCH=amd64 go build -o ./bin/adztbotv2-amd64 ./main.go
env GOOS=linux GOARCH=arm64 go build -o ./bin/adztbotv2-arm64 ./main.go

cd ./bin/

sha256sum adztbotv2-amd64 > adztbotv2-amd64_sha256sum.txt
sha256sum adztbotv2-arm64 > adztbotv2-arm64_sha256sum.txt

zip adztbotv2-amd64 adztbotv2-amd64 adztbotv2-amd64_sha256sum.txt
zip adztbotv2-arm64 adztbotv2-arm64 adztbotv2-arm64_sha256sum.txt