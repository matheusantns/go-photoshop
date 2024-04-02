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

var data psd.InputData

type Steps struct {
	Steps  []int
	Active int
}

var steps = Steps{
	Steps:  []int{1, 2, 3},
	Active: 1,
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "views/static")
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", steps)
	})

	e.POST("/psd/get-fields", func(c echo.Context) error {
		request, err := c.FormParams()
		if err != nil {
			return fmt.Errorf("deu erro - %w", err)
		}

		data = psd.InputData{
			PSExecutableFilePath: request["PSExecutableFilePath"][0],
			ExportDir:            request["ExportDir"][0],
			PSDTemplate:          request["PSDTemplate"][0],
			PrefixNameForFile:    request["PrefixNameForFile"][0],
		}

		data.GetCheckboxValues(request["ExportTypes"])

		fmt.Println(data)

		return c.Render(200, "index", nil)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
