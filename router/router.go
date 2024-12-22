package router

import (
	"database/sql"
	"employees-system/controllers"
	"employees-system/models"
	"employees-system/response"
	"employees-system/session_service"
	"employees-system/utils"
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

var sessionService = session_service.NewSession()

func GetRouter(dbInstance *sql.DB) *chi.Mux {
	const apiVersion = "/api/v1"

	employeeController := controllers.Employee{
		EmployeeService: models.EmployeeService{
			DB: dbInstance,
		},
	}

	roleService := models.RoleService{
		DB: dbInstance,
	}

	technologyService := models.TechnologyService{
		DB: dbInstance,
	}

	projectService := models.ProjectService{
		DB: dbInstance,
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

	funcMap := template.FuncMap{
		"hasRole": func(roleID int, employeeRoles []models.Role) bool {
			for _, r := range employeeRoles {
				if r.ID == roleID {
					return true
				}
			}
			return false
		},
		"hasTechnology": func(techID int, employeeTechs []models.Technology) bool {
			for _, t := range employeeTechs {
				if t.ID == techID {
					return true
				}
			}
			return false
		},
		"hasProject": func(projectID int, employeeProjects []models.Project) bool {
			for _, p := range employeeProjects {
				if p.ID == projectID {
					return true
				}
			}
			return false
		},
	}

	// Admin routes
	adminController := controllers.Admin{
		AdminService: models.AdminService{
			DB: dbInstance,
		},
		SessionService: *sessionService,
	}

	r.Get("/admin/login", func(w http.ResponseWriter, r *http.Request) {
		adminController.RenderLogin(w, r)
	})

	r.Post("/admin/login", func(w http.ResponseWriter, r *http.Request) {
		adminController.Login(w, r)
	})

	r.Get("/admin/test", func(w http.ResponseWriter, r *http.Request) {
		adminController.TestSession(w, r)
	})

	r.Get("/admin/logout", func(w http.ResponseWriter, r *http.Request) {
		adminController.Logout(w, r)
	})

	// Protect routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)

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

			data := utils.TemplateData{
				Employees: employees,
				URL:       utils.URL,
			}

			err = tmpl.ExecuteTemplate(w, "layout-base", data)
			if err != nil {
				fmt.Println(err)
			}
		})

		r.Get(apiVersion+"/employees", func(w http.ResponseWriter, r *http.Request) {
			employeeController.GetAll(w, r)
		})

		r.Get(apiVersion+"/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
			employeeController.GetByID(w, r)
		})

		r.Get("/admin/employees/create", func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := template.ParseFiles("templates/layout-base.html", "templates/employee-create.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			roles, err := roleService.GetAll()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			technologies, err := technologyService.GetAll()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			projectService, err := projectService.GetAll(-1, "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data := utils.TemplateData{
				URL:          utils.URL,
				Roles:        roles,
				Technologies: technologies,
				Projects:     projectService,
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
			tmpl, err := template.New("layout-base").Funcs(funcMap).ParseFiles("templates/layout-base.html", "templates/employee-detail.html")
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

			roles, err := roleService.GetAll()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			technologies, err := technologyService.GetAll()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			projects, err := projectService.GetAll(-1, "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data := utils.TemplateData{
				Employee:     employee,
				Roles:        roles,
				Technologies: technologies,
				Projects:     projects,
				URL:          utils.URL,
			}

			err = tmpl.ExecuteTemplate(w, "layout-base", data)
			if err != nil {
				fmt.Println(err)
			}
		})

		r.Put(apiVersion+"/admin/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
			employeeController.UpdateByID(w, r)
		})

		// Project routes
		projectController := controllers.Project{
			ProjectService: models.ProjectService{
				DB: dbInstance,
			},
		}

		r.Get("/admin/projects", func(w http.ResponseWriter, r *http.Request) {
			projectController.RenderList(w, r)
		})

		r.Get("/admin/projects/create", func(w http.ResponseWriter, r *http.Request) {
			projectController.RenderCreate(w, r)
		})

		r.Post(apiVersion+"/admin/projects/create", func(w http.ResponseWriter, r *http.Request) {
			projectController.Create(w, r)
		})

		r.Get("/admin/projects/{id}", func(w http.ResponseWriter, r *http.Request) {
			projectController.RenderDetailEdit(w, r)
		})

		r.Put(apiVersion+"/admin/projects/{id}", func(w http.ResponseWriter, r *http.Request) {
			projectController.UpdateByID(w, r)
		})

	})

	return r
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		valid := sessionService.CheckLogin(r)

		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
