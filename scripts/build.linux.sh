#!/usr/bin/env bash
cd ..
go build -o demo-rest-api-server -a -ldflags "-extldflags '-static' -X 'main.GitCommit=$(git rev-list -1 HEAD)' -X 'main.GitTag=$(git describe --tags --abbrev=0)'" .
chmod +x ./demo-rest-api-server
echo done.