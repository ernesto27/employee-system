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

	"github.com/go-chi/chi"
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

	nameQuery := r.URL.Query().Get("search")

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

var tplRoleCreate = template.Must(template.ParseFiles("templates/layout-base.html", "templates/roles/role-create.html"))

func (rol *Role) RenderCreate(w http.ResponseWriter, r *http.Request) {

	templateData := utils.TemplateData{
		URL: os.Getenv("URL"),
	}

	err := tplRoleCreate.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}

func (rol *Role) Create(w http.ResponseWriter, r *http.Request) {
	var newRole models.Role
	jsonString := r.FormValue("jsonBody")
	err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&newRole)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	lastID, err := rol.RoleService.Create(newRole)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	newRole.ID = lastID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newRole)
}

var tplRoleEdit = template.Must(template.ParseFiles("templates/layout-base.html", "templates/roles/role-detail-edit.html"))

func (rol *Role) RenderDetailEdit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	role, err := rol.RoleService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := utils.TemplateData{
		URL:  os.Getenv("URL"),
		Role: role,
	}

	err = tplRoleEdit.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}

func (rol *Role) UpdateByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var updatedRole models.Role
	jsonString := r.FormValue("jsonBody")
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&updatedRole)
	if err != nil {
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = rol.RoleService.UpdateByID(id, updatedRole)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRole)
}
