#!/bin/bash

# Nome do executável
EXE_NAME="x-tool-backup"

# Compilar
echo "Compilando o programa..."
go build -o $EXE_NAME main.go
if [ $? -ne 0 ]; then
    echo "Erro ao compilar"
    exit 1
fi

# Rodar
echo "Executando $EXE_NAME..."
./$EXE_NAME
