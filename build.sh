#!/bin/bash

# cross compile...
# windows
GOOS=windows GOARCH=amd64 go build -o build/vmix-mcp.exe cmd/main.go
# mac
GOOS=darwin GOARCH=amd64 go build -o build/vmix-mcp cmd/main.go
