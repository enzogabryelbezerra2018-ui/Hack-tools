# X-Tool Backup

Ferramenta em **Go** para backup e compactação de arquivos, com logs coloridos estilo terminal/Go e geração de imagem do log.

## Funcionalidades

- Compacta todos os arquivos de uma pasta (ex: `x-tool`) em `data.zip`.
- Mostra logs coloridos em tempo real:
  - `[INFO]` → amarelo → em andamento
  - `[OK]` → verde → concluído
  - `[ERR]` → vermelho → erro
- Mantém histórico limitado de logs para simular scroll.
- Gera uma **imagem PNG** (`log.png`) com todos os logs ao final.
- Script para compilar e executar automaticamente (`run.sh` ou `run.bat`).

## Requisitos

- Go 1.21+  
- Linux/macOS/Windows  
- Fonte monoespaçada instalada (para gerar a imagem do log)

## Uso

1. Coloque os arquivos na pasta `x-tool`.  
2. Compile e execute:

```bash
# Linux/macOS
./run.sh

# Windows
run.bat
