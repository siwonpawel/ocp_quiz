package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "embed"
)

//go:embed static/*
var staticContent embed.FS

//go:embed templates/*
var content embed.FS

const (
	HOME_PAGE     = "homePage"
	QUESTION_PAGE = "questionPage"
)

var templates map[string]*template.Template = loadTemplates()

func main() {

	webPort := os.Getenv("WEB_PORT")
	log.Printf("Starting server at localhost:%s", webPort)
	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", webPort),
		Handler: handlers(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}



func loadTemplates() map[string]*template.Template {
	return map[string]*template.Template{
		HOME_PAGE:     template.Must(template.ParseFS(content, "templates/layout.html", "templates/header.html", "templates/mainPage.html", "templates/footer.html")),
		QUESTION_PAGE: template.Must(template.ParseFS(content, "templates/layout.html", "templates/header.html", "templates/questionPage.html", "templates/footer.html")),
	}
}
