package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type ExportType int

const (
	PSD ExportType = iota
	PNG
	JPG
)

type Layer struct {
	Name string
	Type string
}

type FinalData struct {
	PSDTemplate          string
	PSExecutableFilePath string
	ExportDir            string
	ExportTypes          []ExportType
	PrefixNameForFile    string
	Layers               []Layer
	Data                 []map[string]string
}

func (data *FinalData) GetCheckboxValues(str []string) {
	var exportTypes []ExportType

	for _, value := range str {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		exportTypes = append(exportTypes, ExportType(intValue))
	}

	data.ExportTypes = exportTypes
}

func CreateLayers(names []string, layerType string) []Layer {
	var layers []Layer
	for _, name := range names {
		layer := Layer{
			Name: name,
			Type: layerType,
		}
		layers = append(layers, layer)
	}
	return layers
}

func RunPhotoshop(PSExecutableFilePath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("falha ao obter o caminho do executável: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	parentDir := filepath.Clean(filepath.Join(exeDir, ".."))
	scriptPath := filepath.Join(parentDir, "js", "ps_script.js")

	log.Println(parentDir)

	cmd := exec.Command(PSExecutableFilePath, "-r", scriptPath)

	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error - %w", err)
	}

	return nil
}
