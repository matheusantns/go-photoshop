package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	psd "github.com/matheusantns/go-photoshop/internal"
)

type Templates struct {
	templates *template.Template
}

type TextLayer struct {
	TextLayer []psd.Layer
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	t := &Templates{
		templates: template.New("").Funcs(template.FuncMap{
			"ToUpper": strings.ToUpper,
		}),
	}

	t.templates = template.Must(t.templates.ParseGlob("views/*.html"))

	return t
}

func createLayers(names []string, layerType string) []Layer {
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

type Steps struct {
	Steps  []int
	Active int
}

type Error struct {
	Text  string
	Valid bool
}

type Layer struct {
	Name string
	Type string
}

type SecondForm struct {
	Layers []Layer
}

type ThirdForm struct {
	Fields []string
}

type PageData struct {
	Title              string
	FieldToValidate    string
	Steps              Steps
	FirstForm          psd.InputData
	SecondForm         SecondForm
	ThirdForm          ThirdForm
	Error              Error
	AvailableTextLayer []psd.Layer
	BackButtonRoute    string
	ForwardButtonRoute string
}

var pageData PageData

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "views/static")
	e.Renderer = newTemplate()

	pageData.Steps = Steps{
		Steps:  []int{1, 2, 3, 4},
		Active: 1,
	}

	e.GET("/", func(c echo.Context) error {
		pageData.Title = "Preencha as informações a seguir"
		pageData.Steps.Active = 1
		pageData.FieldToValidate = "ExportTypes"
		return c.Render(200, "index", pageData)
	})

	e.GET("/insert-input-image-layer", func(c echo.Context) error {
		return c.Render(200, "input-image-layer.html", nil)
	})

	e.DELETE("/delete-input-image-layer", func(c echo.Context) error {
		return c.NoContent(200)
	})

	e.POST("/step-one", func(c echo.Context) error {
		request, err := c.FormParams()
		if err != nil {
			return fmt.Errorf("deu erro - %w", err)
		}

		pageData.FirstForm = psd.InputData{
			PSExecutableFilePath: request["PSExecutableFilePath"][0],
			ExportDir:            request["ExportDir"][0],
			PSDTemplate:          request["PSDTemplate"][0],
			PrefixNameForFile:    request["PrefixNameForFile"][0],
		}
		pageData.FirstForm.GetCheckboxValues(request["ExportTypes"])

		pageData.AvailableTextLayer, err = psd.HandlePSD(pageData.FirstForm.PSDTemplate)
		if err != nil {
			fmt.Println("deu erro", err)
			pageData.Error = Error{
				Text:  "Arquivo modelo inválido",
				Valid: true,
			}
			return c.Render(200, "form-step-one.html", pageData)
		}

		pageData.Steps.Active = 2
		pageData.Title = "Quais são suas variáveis?"
		pageData.Error = Error{
			Text:  "",
			Valid: false,
		}
		pageData.BackButtonRoute = "/"
		pageData.ForwardButtonRoute = "/step-two"
		pageData.FieldToValidate = "TextLayer"

		return c.Render(200, "form-step-two.html", pageData)
	})

	e.GET("/step-two", func(c echo.Context) error {
		pageData.Steps.Active = 2
		pageData.Title = "Quais são suas variáveis?"
		return c.Render(200, "form-step-two.html", pageData)
	})

	e.POST("/step-two", func(c echo.Context) error {
		request, err := c.FormParams()
		if err != nil {
			return fmt.Errorf("deu erro - %w", err)
		}

		imageLayers := createLayers(request["ImageLayer"], "Image")
		textLayers := createLayers(request["TextLayer"], "Text")
		pageData.BackButtonRoute = "/"
		pageData.ForwardButtonRoute = "/step-two"

		pageData.SecondForm.Layers = append(imageLayers, textLayers...)
		pageData.Steps.Active = 3
		pageData.Title = "Insira o CSV com seus dados"

		return c.Render(200, "form-step-three.html", pageData)
	})

	e.POST("/step-three", func(c echo.Context) error {
		pageData.Steps.Active = 4
		pageData.Title = "Atribua os campos as suas variáveis"
		file, err := c.FormFile("source-csv")
		if err != nil {
			return fmt.Errorf("deu erro - %w", err)
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		csvReader := csv.NewReader(src)
		csvReader.Comma = ';'
		pageData.ThirdForm.Fields, err = csvReader.Read()
		if err != nil {
			log.Fatal(err)
		}

		return c.Render(200, "form-step-four.html", pageData)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
