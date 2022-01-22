@echo off
protoc --go_out=plugins=grpc:. ./*.proto
pause