package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/akshay110495/scanner/loggers"
	"github.com/akshay110495/scanner/service"
)

var (
	logger = loggers.GetLogger()
)

func Home(response http.ResponseWriter, request *http.Request) {
	lp := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		loggers.WriteLog(logger, err)
		return
	}
	tmpl.Execute(response, nil)
}

func Scan(response http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		loggers.WriteLog(logger, err)
		return
	}
	web_url := request.FormValue("web-url")
	main_service := service.GetMainService(web_url)
	links := main_service.GetLinks()
	injectable_links := main_service.GetInjectableLinks()
	response.Header().Set("Content-Type", "application/json")
	data_links := make(map[string][]string)
	data_links["links"] = links
	for url, params := range injectable_links {
		data_links[url] = params
	}
	json.NewEncoder(response).Encode(data_links)
}
