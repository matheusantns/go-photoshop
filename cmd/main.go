package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

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
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Steps struct {
	Steps  []int
	Active int
}

type Error struct {
	Text  string
	Valid bool
}

type SecondForm struct {
	SelectedTextLayers []string
	ImageLayers        []string
}

type PageData struct {
	Title              string
	Steps              Steps
	FirstForm          psd.InputData
	SecondForm         SecondForm
	Error              Error
	AvailableTextLayer []psd.Layer
}

var pageData PageData

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "views/static")
	e.Renderer = newTemplate()

	pageData.Steps = Steps{
		Steps:  []int{1, 2, 3},
		Active: 1,
	}

	e.GET("/", func(c echo.Context) error {
		pageData.Title = "Preencha as informações a seguir"
		pageData.Steps.Active = 1
		return c.Render(200, "index", pageData)
	})

	e.GET("/insert-input-image-layer", func(c echo.Context) error {
		return c.Render(200, "input-image-layer.html", nil)
	})

	e.DELETE("/delete-input-image-layer", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.POST("/step-two", func(c echo.Context) error {
		request, err := c.FormParams()
		if err != nil {
			return fmt.Errorf("deu erro - %w", err)
		}

		pageData.SecondForm.ImageLayers = request["ImageLayer"]
		pageData.SecondForm.SelectedTextLayers = request["TextLayer"]
		fmt.Println(pageData.SecondForm)

		return c.NoContent(http.StatusOK)
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

		return c.Render(200, "form-step-two.html", pageData)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
