@echo off

go build -o disty.exe ./src
cd test
..\\disty.exe init
cd ..