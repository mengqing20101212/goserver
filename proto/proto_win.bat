@echo off
 protoc --proto_path=. --go_out=..\  *.proto
 protoc --proto_path=. --go_out=..\  table\*.proto
set "folderPath=..\protobuf\protobufMsg"
 xcopy  "%folderPath%\*"   ..\src\protobufMsg /E /I

@REM ProtoParser

EM Check if the folder exists
if exist "%folderPath%" (
    echo Deleting folder: %folderPath%
    rmdir /s /q "%folderPath%"
    echo Folder deleted.
) else (
    echo Folder does not exist: %folderPath%
)

::@REM TableParser D:\WORK\me\goserver\proto\table
echo %cd%
set "tableDir=%cd%\..\table"
cd /d "%tableDir%"
dir
@REM for /r %%i in (*Table.go) do (
@REM     echo Formatting file: %%i
@REM     go fmt "%%i"
@REM )