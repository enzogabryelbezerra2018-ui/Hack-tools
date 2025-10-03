package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
)

// Função de log com cores estilo Go
func logInfo(msg string) {
	fmt.Println(ColorYellow + "[INFO] " + msg + ColorReset)
}

func logOK(msg string) {
	fmt.Println(ColorGreen + "[OK]   " + msg + ColorReset)
}

func logErr(msg string) {
	fmt.Println(ColorRed + "[ERR]  " + msg + ColorReset)
}

// Compacta a pasta sourceDir em zipFile
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

func main() {
	sourceDir := "./x-tool" // pasta de origem
	outputZip := "./data.zip"

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		logErr(fmt.Sprintf("Pasta de origem não existe: %s", sourceDir))
		return
	}

	if err := compressToZip(sourceDir, outputZip); err != nil {
		logErr("Erro durante a execução")
	} else {
		logOK(fmt.Sprintf("Todos os arquivos salvos em %s", outputZip))
	}
}
