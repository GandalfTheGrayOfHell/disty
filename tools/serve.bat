@echo off

go build -o disty.exe ./src
echo Running
disty.exe -port=3000 -dir=./tmp serve