package main

import (
	"fmt"
	"html/template"
	"io"

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

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	count := Count{Count: 0}
	layers := Layers{Layers: make([]psd.Layer, 0)}
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		if len(layers.Layers) == 0 {
			fmt.Println("Layers vazia")
		}

		return c.Render(200, "index", struct {
			Count  Count
			Layers Layers
		}{count, layers})
	})

	e.POST("/psd/get-fields", func(c echo.Context) error {
		psdPath := c.FormValue("psd")
		layersResult, err := psd.HandlePSD(psdPath)
		if err != nil {
			return fmt.Errorf("deu erro - %w", err)
		}

		layers.Layers = layersResult

		if len(layers.Layers) == 0 {
			fmt.Println("Layers vazia")
		}

		return c.Render(200, "index", layers)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
