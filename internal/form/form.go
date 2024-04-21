package form

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"

	psd "github.com/matheusantns/go-photoshop/internal/psd"
	utils "github.com/matheusantns/go-photoshop/internal/utils"
)

type Steps struct {
	Steps  []int
	Active int
}

type Error struct {
	Text  string
	Valid bool
}

type ThirdForm struct {
	Fields []string
	Data   [][]string
}

type PageData struct {
	Title              string
	Steps              Steps
	FirstForm          psd.PSDTemplate
	ThirdForm          ThirdForm
	Error              Error
	AvailableTextLayer []psd.Layer
	FinalData          utils.FinalData
}

func (data *PageData) PopulateThirdForm(src multipart.File) error {
	csvReader := csv.NewReader(src)
	csvReader.Comma = ';'

	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("empty CSV file")
	}

	data.ThirdForm.Data = records
	data.ThirdForm.Fields = data.ThirdForm.Data[0]

	return nil
}

func ReadCSV(file *multipart.FileHeader) (multipart.File, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	return src, nil
}
