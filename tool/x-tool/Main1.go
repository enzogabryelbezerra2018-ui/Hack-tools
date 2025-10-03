package main

import (
	"archive/zip"
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"

	MaxLogLines = 30 // linhas visíveis no terminal
)

var logLines []string // histórico completo de logs

// Funções de log colorido
func logInfo(msg string) {
	appendLog("[INFO] " + msg)
	fmt.Println(ColorYellow + "[INFO] " + msg + ColorReset)
}

func logOK(msg string) {
	appendLog("[OK]   " + msg)
	fmt.Println(ColorGreen + "[OK]   " + msg + ColorReset)
}

func logErr(msg string) {
	appendLog("[ERR]  " + msg)
	fmt.Println(ColorRed + "[ERR]  " + msg + ColorReset)
}

// Adiciona log ao histórico e mantém tamanho limitado (para simular scroll)
func appendLog(line string) {
	logLines = append(logLines, line)
	if len(logLines) > MaxLogLines {
		logLines = logLines[len(logLines)-MaxLogLines:]
	}
}

// Compacta pasta sourceDir para zipFile
func compressToZip(sourceDir, zipFile string) error {
	logInfo("Iniciando compactação...")
	zipOut, err := os.Create(zipFile)
	if err != nil {
		logErr(fmt.Sprintf("Erro ao criar zip: %v", err))
		return err
	}
	defer zipOut.Close()

	archive := zip.NewWriter(zipOut)
	defer archive.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logErr(fmt.Sprintf("Erro acessando %s: %v", path, err))
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			logErr(fmt.Sprintf("Erro calculando caminho relativo: %v", err))
			return err
		}

		logInfo(fmt.Sprintf("Compactando: %s", relPath))
		time.Sleep(150 * time.Millisecond) // simula delay de processamento

		file, err := os.Open(path)
		if err != nil {
			logErr(fmt.Sprintf("Erro abrindo arquivo: %v", err))
			return err
		}
		defer file.Close()

		f, err := archive.Create(relPath)
		if err != nil {
			logErr(fmt.Sprintf("Erro adicionando arquivo ao zip: %v", err))
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			logErr(fmt.Sprintf("Erro copiando conteúdo: %v", err))
			return err
		}

		logOK(fmt.Sprintf("%s concluído", relPath))
		return nil
	})

	if err != nil {
		logErr(fmt.Sprintf("Falha durante a compactação: %v", err))
		return err
	}

	logOK(fmt.Sprintf("Compactação finalizada: %s", zipFile))
	return nil
}

// Gera uma imagem PNG com todos os logs
func generateLogImage(output string) error {
	const W, H, padding = 800, 600, 10
	dc := gg.NewContext(W, H)
	dc.SetColor(color.Black)
	dc.Clear()

	dc.SetColor(color.White)
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf", 16); err != nil {
		fmt.Println("Erro carregando fonte:", err)
	}

	y := padding
	for _, line := range logLines {
		col := color.White
		if len(line) >= 6 {
			switch line[:5] {
			case "[INF":
				col = color.RGBA{255, 255, 0, 255}
			case "[OK]":
				col = color.RGBA{0, 255, 0, 255}
			case "[ERR]":
				col = color.RGBA{255, 0, 0, 255}
			}
		}
		dc.SetColor(col)
		dc.DrawString(line, float64(padding), float64(y))
		y += 20
		if y > H-padding {
			break
		}
	}

	return dc.SavePNG(output)
}

func main() {
	sourceDir := "./x-tool"
	outputZip := "./data.zip"
	outputImage := "./log.png"

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		logErr(fmt.Sprintf("Pasta de origem não existe: %s", sourceDir))
		return
	}

	err := compressToZip(sourceDir, outputZip)
	if err != nil {
		logErr("Erro durante compactação")
	} else {
		logOK(fmt.Sprintf("Todos os arquivos salvos em %s", outputZip))
	}

	time.Sleep(time.Second)
	err = generateLogImage(outputImage)
	if err != nil {
		logErr(fmt.Sprintf("Erro gerando imagem do log: %v", err))
	} else {
		logOK(fmt.Sprintf("Imagem do log criada: %s", outputImage))
	}

	logInfo("Processo finalizado.")
}
