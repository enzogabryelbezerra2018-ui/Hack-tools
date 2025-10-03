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


3.	Aguarde a compactação terminar.
	4.	Confira:
	•	Backup: data.zip
	•	Log visual: log.png
---

## 2️⃣ `USB-logic.go`

Esse arquivo vai conter **lógica para detectar se um dispositivo USB está conectado** (para simular o alerta do seu programa).  
Observação: Go puro não acessa diretamente USB no Android ou Windows sem libs externas. Aqui vai uma **simulação básica para Linux/macOS**:

```go
package main

import (
	"fmt"
	"io/ioutil"
	"time"
)

// Pasta de exemplo onde dispositivos montados aparecem (Linux)
const usbMountPath = "/media" // no Windows, poderia ser "E:\\"

func checkUSBConnected() bool {
	files, err := ioutil.ReadDir(usbMountPath)
	if err != nil {
		return false
	}
	return len(files) > 0
}

func waitForUSB() {
	fmt.Println("[INFO] Aguardando dispositivo USB...")
	for {
		if checkUSBConnected() {
			fmt.Println("[OK] Dispositivo USB detectado!")
			return
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	waitForUSB()
	fmt.Println("[INFO] Pode iniciar a ferramenta de backup")
}
