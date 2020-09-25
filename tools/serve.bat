@echo off

go build -o disty.exe ./src
echo Running
disty.exe serve -port=3000 -dir=./tmp