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
	generator := getTemplateGenerator(uiTitle, uiDescription, template, services)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", redirectToServiceHandler(services))
	r.Get("/{name}", serveUiHandler(services, generator))
	r.Get("/file/{file}", readFileHandler(uiFilePath, services))

	if err = http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		panic(err)
	}
}

func redirectToServiceHandler(services ServiceList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := services[0]
		http.Redirect(w, r, service.DocUrl, http.StatusFound)
	}
}

func serveUiHandler(services ServiceList, generator templateGenerator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		index := services.Find(chi.URLParam(request, "name"), func(service *Service) string {
			return service.Name
		})
		if index == -1 {
			http.NotFound(writer, request)
			return
		}
		_ = generator(writer, services[index])
	}
}

func readFileHandler(filePath string, services ServiceList) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		index := services.Find(chi.URLParam(request, "file"), func(service *Service) string {
			return service.File
		})
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
