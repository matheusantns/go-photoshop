package main

import (
	"fmt"
	"html/template"
	"io"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	psd "github.com/matheusantns/go-photoshop/internal"
)

type Templates struct {
	templates *template.Template
}

type Count struct {
	Count int
}

type Layers struct {
	Layers []psd.Layer
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

// layersResult, err := psd.HandlePSD(psdModel)
// 		if err != nil {
// 			return fmt.Errorf("deu erro - %w", err)
// 		}

// 		layers := Layers{Layers: layersResult}

func getCheckboxValues(str []string) []psd.ExportType {
	var exportTypes []psd.ExportType

	for _, value := range str {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		exportTypes = append(exportTypes, psd.ExportType(intValue))
	}

	return exportTypes
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "views/static")
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.POST("/psd/get-fields", func(c echo.Context) error {
		request, err := c.FormParams()
		if err != nil {
			return fmt.Errorf("deu erro - %w", err)
		}

		exportTypes := getCheckboxValues(request["ExportTypes"])

		data := psd.InputData{
			PSExecutableFilePath: request["PSExecutableFilePath"][0],
			ExportDir:            request["ExportDir"][0],
			ExportTypes:          exportTypes,
			PSDTemplate:          request["PSDTemplate"][0],
			PrefixNameForFile:    request["PrefixNameForFile"][0],
		}

		fmt.Println(data)

		return c.Render(200, "index", nil)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
