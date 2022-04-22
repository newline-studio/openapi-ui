package main

import (
	_ "embed"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	uiTitle := os.Getenv("UI_TITLE")
	uiDescription := os.Getenv("UI_DESCRIPTION")
	uiFilePath := os.Getenv("UI_FILE_PATH")
	uiUrl := os.Getenv("UI_URL")
	services, err := getServicesFromString(os.Getenv("UI_SERVICES"), uiUrl)
	if err != nil {
		panic(err)
	}

	template, err := getTemplate()
	if err != nil {
		panic(err)
	}
	templateResult, err := prepareTemplateString(uiTitle, uiDescription, template, services)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", serveUiHandler(templateResult))
	r.Get("/file/{file}", readFileHandler(uiFilePath, services))

	if err = http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		panic(err)
	}
}

func serveUiHandler(template string) http.HandlerFunc {
	tplBytes := []byte(template)
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(tplBytes)
	}
}

func readFileHandler(filePath string, services ServiceList) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		index := services.Find(chi.URLParam(request, "file"))
		if index == -1 {
			http.NotFound(writer, request)
			return
		}
		file, err := os.ReadFile(filePath + "/" + services[index].File)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(file)
	}
}
