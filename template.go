package main

import (
	"embed"
	"html/template"
	"io"
)

type templateGenerator func(io.Writer, Service) error

//go:embed ui.gohtml
var uiFs embed.FS

type templateData struct {
	Title       string
	Description string
	File        string
	Download    string
	Services    ServiceList
	Service     Service
}

func getTemplate() (*template.Template, error) {
	data, err := uiFs.ReadFile("ui.gohtml")
	if err != nil {
		panic(err)
	}
	return template.New("ui").Parse(string(data))
}

func getTemplateGenerator(title, description string, tpl *template.Template, services ServiceList) templateGenerator {
	return func(writer io.Writer, service Service) error {
		return tpl.Execute(writer, templateData{
			Title:       title,
			Description: description,
			File:        service.FileUrl,
			Download:    service.DownloadFileUrl,
			Services:    services,
			Service:     service,
		})
	}
}
