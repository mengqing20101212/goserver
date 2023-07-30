@echo off
 ..\tool\protoc-23.4-win64\bin\protoc.exe --proto_path=. --go_out=..\  *.proto
 ..\tool\protoc-23.4-win64\bin\protoc.exe --proto_path=. --go_out=..\  table\*.proto
ProtoParser
TableParser G:\WORK\me\goserver\proto\table
echo %cd%
set "tableDir=%cd%\..\table"
cd /d "%tableDir%"
dir
for /r %%i in (*Table.go) do (
    echo Formatting file: %%i
    go fmt "%%i"
)