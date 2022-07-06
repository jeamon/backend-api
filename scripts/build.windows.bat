@echo off
cd ..
git rev-list -1 HEAD > tempFile && set /P gitCommit=<tempFile
git describe --tags --abbrev=0 > tempFile && set /P gitTag=<tempFile
del tempFile
go build -o demo-rest-api-server.exe -a -ldflags "-extldflags '-static' -X 'main.GitCommit=%gitCommit%' -X 'main.GitTag=%gitTag%'" .
echo done.