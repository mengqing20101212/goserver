@echo off
set output_dir="G:\WORK\me\goserver\common\proto"
 G:\WORK\me\goserver\tool\protoc-23.4-win64\bin\protoc.exe --proto_path=. --go_out=%output_dir% *.proto