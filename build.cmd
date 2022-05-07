@echo off
set GOOS=linux
set CGO_ENABLED=0
go build -trimpath -o quic-wget

set GOOS=windows
go build -trimpath -o quic-wget.exe
