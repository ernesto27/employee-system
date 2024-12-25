package controllers

import (
	"employees-system/models"
	"employees-system/session_service"
	"employees-system/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Admin struct {
	AdminService   models.AdminService
	SessionService session_service.Session
}

func (admin *Admin) RenderLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := utils.TemplateData{
		URL: os.Getenv("URL"),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (admin *Admin) Login(w http.ResponseWriter, r *http.Request) {
	var adminRequest models.Admin
	jsonString := r.FormValue("jsonBody")
	err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&adminRequest)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	_, err = admin.AdminService.Login(adminRequest.Email, adminRequest.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = admin.SessionService.CreateLogin(w, r)

	if err != nil {
		http.Error(w, "Unable to save session", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("login success"))
}

func (admin *Admin) TestSession(w http.ResponseWriter, r *http.Request) {
	v := admin.SessionService.CheckLogin(r)
	w.Write([]byte(fmt.Sprintf("authenticated: %v", v)))
}

func (admin *Admin) Logout(w http.ResponseWriter, r *http.Request) {
	err := admin.SessionService.RemoveLogin(w, r)
	if err != nil {
		http.Error(w, "Unable to remove session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/login", http.StatusOK)
}
