package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
)

//go:embed ui.gohtml
var uiFs embed.FS

type templateData struct {
	Title       string
	Description string
	Services    template.JS
}

func getTemplate() (*template.Template, error) {
	data, err := uiFs.ReadFile("ui.gohtml")
	if err != nil {
		panic(err)
	}
	return template.New("ui").Parse(string(data))
}

func prepareTemplateString(title, description string, tpl *template.Template, services []Service) (string, error) {
	jsonServices, err := json.Marshal(services)
	if err != nil {
		return "", err
	}
	serviceString := string(jsonServices)
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, templateData{
		Title:       title,
		Description: description,
		Services:    template.JS(serviceString),
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
