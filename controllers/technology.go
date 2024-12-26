package controllers

import (
	"employees-system/models"
	"employees-system/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Technology struct {
	TechnologyService models.TechnologyService
}

var tmplTechnologyList = template.Must(template.ParseFiles("templates/layout-base.html", "templates/technologies/technology-list.html"))

func (technology *Technology) RenderList(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")
	pageNum := 1

	if page != "" {
		var err error
		pageNum, err = strconv.Atoi(page)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	technologies, err := technology.TechnologyService.GetAll(pageNum, search)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := utils.TemplateData{
		URL:          os.Getenv("URL"),
		Technologies: technologies,
	}

	err = tmplTechnologyList.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var tmplTechnologyCreate = template.Must(template.ParseFiles("templates/layout-base.html", "templates/technologies/technology-create.html"))

func (technology *Technology) RenderCreate(w http.ResponseWriter, r *http.Request) {
	templateData := utils.TemplateData{
		URL: os.Getenv("URL"),
	}
	err := tmplTechnologyCreate.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (technology *Technology) Create(w http.ResponseWriter, r *http.Request) {
	var newTechnology models.Technology
	jsonString := r.FormValue("jsonBody")
	err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&newTechnology)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	insertID, err := technology.TechnologyService.Create(newTechnology)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newTechnology.ID = insertID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTechnology)
}
