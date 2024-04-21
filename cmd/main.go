package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	form "github.com/matheusantns/go-photoshop/internal/form"
	psd "github.com/matheusantns/go-photoshop/internal/psd"
	utils "github.com/matheusantns/go-photoshop/internal/utils"
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

var pageData form.PageData

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "views/static")
	e.Renderer = newTemplate()

	pageData.Steps = form.Steps{
		Steps:  []int{1, 2, 3, 4},
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
		return c.NoContent(200)
	})

	e.POST("/step-one", func(c echo.Context) error {
		request, err := c.FormParams()
		if err != nil {
			return fmt.Errorf("error - %w", err)
		}

		pageData.FinalData = utils.FinalData{
			PSExecutableFilePath: request["PSExecutableFilePath"][0],
			ExportDir:            request["ExportDir"][0],
			PrefixNameForFile:    request["PrefixNameForFile"][0],
			PSDTemplate:          request["PSDTemplate"][0],
		}

		pageData.FinalData.GetCheckboxValues(request["ExportTypes"])

		pageData.AvailableTextLayer, err = psd.HandlePSD(pageData.FinalData.PSDTemplate)
		if err != nil {
			fmt.Println("Error", err)
			pageData.Error = form.Error{
				Text:  "Arquivo modelo inválido",
				Valid: true,
			}
			return c.Render(200, "form-step-one.html", pageData)
		}

		pageData.Steps.Active = 2
		pageData.Title = "Quais são suas variáveis?"
		pageData.Error = form.Error{
			Text:  "",
			Valid: false,
		}

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
			return fmt.Errorf("error - %w", err)
		}

		imageLayers := utils.CreateLayers(request["ImageLayer"], "Image")
		textLayers := utils.CreateLayers(request["TextLayer"], "Text")
		pageData.FinalData.Layers = append(imageLayers, textLayers...)

		pageData.Steps.Active = 3
		pageData.Title = "Insira o CSV com seus dados"

		return c.Render(200, "form-step-three.html", pageData)
	})

	e.GET("/step-three", func(c echo.Context) error {
		pageData.Steps.Active = 3
		pageData.Title = "Insira o CSV com seus dados"
		return c.Render(200, "form-step-three.html", pageData)
	})

	e.POST("/step-three", func(c echo.Context) error {
		pageData.Steps.Active = 4
		pageData.Title = "Atribua os campos as suas variáveis"

		file, err := c.FormFile("source-csv")
		if err != nil {
			return fmt.Errorf("error - %w", err)
		}

		fileRead, nil := form.ReadCSV(file)
		if err != nil {
			return fmt.Errorf("error - %w", err)
		}

		pageData.PopulateThirdForm(fileRead)

		return c.Render(200, "form-step-four.html", pageData)
	})

	e.POST("/step-four", func(c echo.Context) error {
		request, err := c.FormParams()
		if err != nil {
			return fmt.Errorf("error - %w", err)
		}

		mapping := make(map[string]string)
		for _, values := range request {
			newKey := values[1]
			newValue := values[0]
			mapping[newKey] = newValue
		}

		results := []map[string]string{}

		for _, line := range pageData.ThirdForm.Data[1:] {
			item := map[string]string{}

			for i, value := range line {
				defaultHeader := pageData.ThirdForm.Data[0][i]
				newHeader, ok := mapping[defaultHeader]
				if !ok {
					continue
				}
				item[newHeader] = value
			}

			results = append(results, item)
		}

		pageData.FinalData.Data = results

		jason, err := json.Marshal(pageData.FinalData)
		if err != nil {
			return fmt.Errorf("error - %w", err)
		}

		err = os.WriteFile("js\\parameters.json", jason, 0644)
		if err != nil {
			return fmt.Errorf("error writing JSON data to file: %w", err)
		}

		utils.RunPhotoshop(pageData.FinalData.PSExecutableFilePath)

		return c.Render(200, "processing.html", nil)
	})

	e.Logger.Fatal(e.Start(":3001"))
}
