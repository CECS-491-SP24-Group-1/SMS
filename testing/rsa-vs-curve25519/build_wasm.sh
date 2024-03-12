#!/bin/bash
BIN_NAME=benchmark.wasm
ENTRYPOINT_NAME=main_wasm.go
GOARCH=wasm GOOS=js go build -ldflags "-s -w" -o $BIN_NAME $ENTRYPOINT_NAME
