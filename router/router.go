package router

import (
	"database/sql"
	"employees-system/controllers"
	"employees-system/models"
	"employees-system/response"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

type TemplateData struct {
	Employees []models.Employee
	URL       string
}

func GetRouter(dbInstance *sql.DB) *chi.Mux {
	const apiVersion = "/api/v1"

	employeeController := controllers.Employee{
		EmployeeService: models.EmployeeService{
			DB: dbInstance,
		},
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(c.Handler)

	// Serve static files from the "public" directory
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	FileServer(r, "/admin/public", filesDir)

	r.Get(apiVersion+"/", func(w http.ResponseWriter, r *http.Request) {
		response.NewWithoutData().WithMessage("employees system api v1").Success(w)
	})

	r.Get(apiVersion+"/employees", func(w http.ResponseWriter, r *http.Request) {
		employeeController.GetAll(w, r)
	})

	r.Get(apiVersion+"/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
		employeeController.GetByID(w, r)
	})

	url := "http://localhost:8080/admin/public"

	r.Get("/admin/employees", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/layout-base.html", "templates/employee-list.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}

		employees, err := employeeController.EmployeeService.GetAll(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			Employees: employees,
			URL:       url,
		}

		err = tmpl.ExecuteTemplate(w, "layout-base", data)
		if err != nil {
			fmt.Println(err)
		}
	})

	r.Get("/admin/employees/create", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/layout-base.html", "templates/employee-create.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			URL: url,
		}

		err = tmpl.ExecuteTemplate(w, "layout-base", data)
		if err != nil {
			fmt.Println(err)
		}
	})

	r.Post(apiVersion+"/admin/employees/create", func(w http.ResponseWriter, r *http.Request) {
		employeeController.Create(w, r)
	})

	r.Get("/admin/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/employee_detail.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid employee ID", http.StatusBadRequest)
			return
		}

		employee, err := employeeController.EmployeeService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, employee)
	})

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, fs.ServeHTTP)
}

func ValidateAdminApiToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getApiToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token")
	}
	token := splitToken[1]
	return token, nil
}
