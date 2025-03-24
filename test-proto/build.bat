@echo off

REM 切換到批次檔所在的目錄
cd /d %~dp0

set out_dir="./dist"

if exist %out_dir% (
    echo remove dir dist
    rmdir /q /s %out_dir%
) 

mkdir %out_dir%

set PATH=%PATH%;./bin/

protoc --proto_path="./demo" --go_out=%out_dir% ^
demo/source.proto