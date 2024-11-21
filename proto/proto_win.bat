@echo off
 protoc --proto_path=. --go_out=..\  *.proto
 protoc --proto_path=. --go_out=..\  table\*.proto
set "folderPath=..\protobuf\protobufMsg"
xcopy  "%folderPath%\*"   ..\src\protobufMsg /E /I
xcopy  "..\table\*"   ..\src\table /E /I
call delFolder "%folderPath%"
call delFolder "..\table"
ProtoParser



TableParser .
echo %cd%
set "tableDir=%cd%\..\src\table"
cd /d "%tableDir%"
dir
 for /r %%i in (*Table.go) do (
     echo Formatting file: %%i
     go fmt "%%i"
 )

 :delFolder
  set "folderPath=%1"
  REM Check if the folder exists
  if exist "%folderPath%" (
      echo Deleting folder: %folderPath%
      rmdir /s /q "%folderPath%"
      echo Folder deleted.
  ) else (
      echo Folder does not exist: %folderPath%
  )
  goto :eof