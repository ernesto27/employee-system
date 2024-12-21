package controllers

import (
	"employees-system/models"
	"employees-system/utils"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Project struct {
	ProjectService models.ProjectService
}

var tmplCreate = template.Must(template.ParseFiles("templates/layout-base.html", "templates/projects/project-list.html"))

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

	projects, err := project.ProjectService.GetAll(page)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := utils.TemplateData{
		URL:      utils.URL,
		Projects: projects,
	}

	err = tmplCreate.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}
