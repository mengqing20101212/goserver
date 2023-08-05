@echo off
ConfigParser
echo %cd%
set "tableDir=%cd%\..\config"
cd /d "%tableDir%"
dir
go fmt
