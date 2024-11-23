@echo off
ConfigParser .
set "tableDir=%cd%\..\src\config"
cd /d "%tableDir%"
echo "CurDir:" %cd%
dir
go fmt
