@echo off
set EXE_NAME=x-tool-backup.exe

echo Compilando programa...
go build -o %EXE_NAME% main.go
if errorlevel 1 (
    echo Erro ao compilar
    exit /b 1
)

echo Executando %EXE_NAME%...
%EXE_NAME%
pause
