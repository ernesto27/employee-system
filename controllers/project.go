package controllers

import (
	"employees-system/models"
	"employees-system/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Project struct {
	ProjectService models.ProjectService
}

var tmplList = template.Must(template.ParseFiles("templates/layout-base.html", "templates/projects/project-list.html"))

func (project *Project) RenderList(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	nameQuery := r.URL.Query().Get("name")

	projects, err := project.ProjectService.GetAll(page, nameQuery)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := utils.TemplateData{
		URL:      utils.URL,
		Projects: projects,
	}

	err = tmplList.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}

var tmplCreate = template.Must(template.ParseFiles("templates/layout-base.html", "templates/projects/project-create.html"))

func (project *Project) RenderCreate(w http.ResponseWriter, r *http.Request) {
	templateData := utils.TemplateData{
		URL: utils.URL,
	}

	err := tmplCreate.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}

func (project *Project) Create(w http.ResponseWriter, r *http.Request) {
	var newProject models.Project
	jsonString := r.FormValue("jsonBody")
	err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&newProject)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = project.ProjectService.Create(newProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProject)

}
