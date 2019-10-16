@echo off
cd /d %~dp0

START "data collector" ^
main.exe
