package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ====== LOGS ======
var logs []string

func logInfo(msg string) {
	line := fmt.Sprintf("[INFO] %s", msg)
	fmt.Println("\033[33m" + line + "\033[0m") // amarelo
	logs = append(logs, line)
}

func logOK(msg string) {
	line := fmt.Sprintf("[OK] %s", msg)
	fmt.Println("\033[32m" + line + "\033[0m") // verde
	logs = append(logs, line)
}

func logErr(msg string) {
	line := fmt.Sprintf("[ERR] %s", msg)
	fmt.Println("\033[31m" + line + "\033[0m") // vermelho
	logs = append(logs, line)
}

// ====== USB DETECT ======
const usbMountPath = "/media" // para Linux/macOS; Windows seria "E:\\"

func checkUSBConnected() bool {
	files, err := ioutil.ReadDir(usbMountPath)
	if err != nil {
		return false
	}
	return len(files) > 0
}

func waitForUSB() {
	logInfo("Aguardando dispositivo USB...")
	for {
		if checkUSBConnected() {
			logOK("Dispositivo USB detectado!")
			return
		}
		time.Sleep(2 * time.Second)
	}
}

// ====== ZIP ======
func addFileToZip(zipWriter *zip.Writer, filename, baseInZip string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Join(baseInZip, filepath.Base(filename))
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

func compressFolderToZip(sourceDir, outputZip string) error {
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			logInfo("Compactando " + path)
			err = addFileToZip(zipWriter, path, filepath.Base(sourceDir))
			if err != nil {
				return err
			}
			logOK("Arquivo adicionado: " + path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// ====== GERAR LOG.PNG ======
func saveLogsAsImage(filename string) error {
	width := 800
	height := 20*len(logs) + 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// Renderiza o texto como quadrados de cor (simples)
	y := 20
	for _, line := range logs {
		var c color.Color = color.Black
		if strings.HasPrefix(line, "[INFO]") {
			c = color.RGBA{255, 200, 0, 255}
		} else if strings.HasPrefix(line, "[OK]") {
			c = color.RGBA{0, 200, 0, 255}
		} else if strings.HasPrefix(line, "[ERR]") {
			c = color.RGBA{200, 0, 0, 255}
		}
		for x := 10; x < width-10; x++ {
			img.Set(x, y, c)
		}
		y += 20
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, img)
}

// ====== MAIN ======
func main() {
	waitForUSB()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Você deseja iniciar o backup? (y/n): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice != "y" && choice != "Y" {
		logErr("Operação cancelada pelo usuário.")
		return
	}

	sourceDir := "x-tool"
	outputZip := "data.zip"

	logInfo("Iniciando backup...")
	err := compressFolderToZip(sourceDir, outputZip)
	if err != nil {
		logErr("Erro: " + err.Error())
		return
	}
	logOK("Backup concluído: " + outputZip)

	logInfo("Gerando log visual...")
	err = saveLogsAsImage("log.png")
	if err != nil {
		logErr("Erro ao gerar log visual: " + err.Error())
		return
	}
	logOK("Log salvo em log.png")
}
