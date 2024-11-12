@echo off
 protoc --proto_path=. --go_out=..\  *.proto
 protoc --proto_path=. --go_out=..\  table\*.proto

 xcopy ..\protobuf\protobufMsg\*   ..\src\protobufMsg /E /I

@REM ProtoParser
@REM TableParser D:\WORK\me\goserver\proto\table
echo %cd%
set "tableDir=%cd%\..\table"
cd /d "%tableDir%"
dir
@REM for /r %%i in (*Table.go) do (
@REM     echo Formatting file: %%i
@REM     go fmt "%%i"
@REM )