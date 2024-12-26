package controllers

import (
	"employees-system/models"
	"employees-system/utils"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
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
