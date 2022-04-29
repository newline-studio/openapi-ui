package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/newline-studio/swagger-ui/resolver"
)

func main() {
	uiTitle := os.Getenv("UI_TITLE")
	uiDescription := os.Getenv("UI_DESCRIPTION")
	uiFilePath := os.Getenv("UI_FILE_PATH")
	uiUrl := os.Getenv("UI_URL")
	uiPath := os.Getenv("UI_PATH")
	services, err := getServicesFromString("UI_SERVICE_", os.Environ(), uiUrl)

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
	r.Get("/files/*", noCacheMiddleware(readFileHandler(uiFilePath, uiPath+"/files")))
	r.Get("/downloads/*", noCacheMiddleware(downloadFileHandler(uiFilePath, uiPath+"/files", uiPath+"/downloads")))

	log.Println("Starting openapi-ui server on port 8080...")
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

func downloadFileHandler(filePath string, prefix string, initial string) http.HandlerFunc {
	pathResolver := func(path string) string {
		return strings.Replace(path, prefix, filePath, 1)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		path := strings.Replace(request.URL.Path, initial, prefix, 1)
		r := resolver.NewResolver(pathResolver)
		result, err := r.Resolve(path)
		if err != nil {
			log.Println("Error resolving path:", err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(result)
	}
}

func readFileHandler(filePath string, prefix string) http.HandlerFunc {
	handler := http.StripPrefix(prefix, http.FileServer(http.Dir(filePath)))
	return func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	}
}

func noCacheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	expires := time.Unix(0, 0).Format(time.RFC1123)
	return func(writer http.ResponseWriter, request *http.Request) {
		headers := writer.Header()
		headers.Add("Expires", expires)
		headers.Add("Cache-Control", "no-cache")
		headers.Add("Pragma", "no-cache")
		headers.Add("X-Accel-Expires", "0")
		next(writer, request)
	}
}
