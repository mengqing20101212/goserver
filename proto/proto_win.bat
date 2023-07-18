@echo off
 ..\tool\protoc-23.4-win64\bin\protoc.exe --proto_path=. --go_out=..\  *.proto
