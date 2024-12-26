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

type Role struct {
	RoleService models.RoleService
}

var tplRolList = template.Must(template.ParseFiles("templates/layout-base.html", "templates/roles/role-list.html"))

func (rol *Role) RenderList(w http.ResponseWriter, r *http.Request) {
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

	roles, err := rol.RoleService.GetAll(page, nameQuery)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := utils.TemplateData{
		URL:   os.Getenv("URL"),
		Roles: roles,
	}

	err = tplRolList.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}
