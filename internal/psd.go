package psd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type ExportType int
type Layer string

const (
	PSD ExportType = iota
	PNG
	JPG
)

type InputData struct {
	PSExecutableFilePath string
	ExportDir            string
	ExportTypes          []ExportType
	PSDTemplate          string
	PrefixNameForFile    string
	// DataToRenderOnPSD    map[string]interface{}
}

func FindPSDHeader(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var data []byte
	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		data = append(data, buffer[:n]...)

		if idx := bytes.Index(data, []byte("</x:xmpmeta>")); idx != -1 {
			start := bytes.Index(data, []byte("<x:xmpmeta"))
			data = data[start : idx+len("</x:xmpmeta>")]
			break
		}

		if err == io.EOF {
			break
		}
	}

	return data, err
}

func FindTextLayers(xmp string) ([]Layer, error) {
	substr := "photoshop:LayerName"
	indices := []int{}
	start := 0
	layers := []Layer{}

	for {
		idx := strings.Index(xmp[start:], substr)
		if idx == -1 {
			break
		}
		indices = append(indices, start+idx)
		start += idx + len(substr)
	}

	j := 1
	for i := 0; i <= len(indices)-1; i += 2 {
		layer := xmp[(indices[i] + len(substr) + 1) : indices[j]-2]
		layers = append(layers, Layer(layer))
		j += 2
	}

	return layers, nil
}

func HandlePSD(filepath string) ([]Layer, error) {
	// filepath := "C:\\Users\\teteu\\OneDrive\\Documentos\\Tarkuss\\Guia_dos_Precos\\00 - Modelo - Patrulha da Pascoa.psd"

	data, err := FindPSDHeader(filepath)
	if err != nil {
		fmt.Println("Erro ao abrir", err)
		return nil, fmt.Errorf("deu erro - %w", err)
	}

	layers, err := FindTextLayers(string(data))
	if err != nil {
		fmt.Println("Erro ao encontrar camadas de texto", err)
		return nil, fmt.Errorf("deu erro - %w", err)
	}

	return layers, nil
}
