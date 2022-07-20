@echo off
protoc --go-grpc_out=. --go_out=. ./*.proto
pause