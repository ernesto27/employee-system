package router

import (
	"database/sql"
	"employees-system/controllers"
	"employees-system/internal/s3"
	"employees-system/models"
	"employees-system/response"
	"employees-system/session_service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

var sessionService = session_service.NewSession()

func GetRouter(dbInstance *sql.DB, myS3 *s3.MyS3, url string) *chi.Mux {
	const apiVersion = "/api/v1"

	roleService := models.RoleService{
		DB: dbInstance,
	}

	technologyService := models.TechnologyService{
		DB: dbInstance,
	}

	projectService := models.ProjectService{
		DB: dbInstance,
	}

	employeeController := controllers.Employee{
		EmployeeService: models.EmployeeService{
			DB: dbInstance,
		},
		RoleService:       roleService,
		TechnologyService: technologyService,
		ProjectService:    projectService,
		S3Service:         *myS3,
		ImageService: models.ImageService{
			DB: dbInstance,
		},
	}

	technologyController := controllers.Technology{
		TechnologyService: models.TechnologyService{
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
			employeeController.RenderList(w, r)
		})

		r.Get(apiVersion+"/employees", func(w http.ResponseWriter, r *http.Request) {
			employeeController.GetAll(w, r)
		})

		r.Get(apiVersion+"/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
			employeeController.GetByID(w, r)
		})

		r.Get("/admin/employees/create", func(w http.ResponseWriter, r *http.Request) {
			employeeController.RenderCreate(w, r)
		})

		r.Post(apiVersion+"/admin/employees/create", func(w http.ResponseWriter, r *http.Request) {
			employeeController.Create(w, r)
		})

		r.Get("/admin/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
			employeeController.RenderDetailEdit(w, r)
		})

		r.Put(apiVersion+"/admin/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
			employeeController.UpdateByID(w, r)
		})

		// Project routes
		projectController := controllers.Project{
			ProjectService: models.ProjectService{
				DB: dbInstance,
			},
			S3Service: *myS3,
			ImageService: models.ImageService{
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

		r.Get(apiVersion+"/admin/images", func(w http.ResponseWriter, r *http.Request) {
			queryPath := r.URL.Query().Get("path")
			if queryPath == "" {
				http.Error(w, "Invalid path", http.StatusBadRequest)
				return
			}

			b, err := myS3.Get(queryPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(b)

		})

		// Technology routes
		r.Get("/admin/technologies", func(w http.ResponseWriter, r *http.Request) {
			technologyController.RenderList(w, r)
		})

		r.Get("/admin/technologies/create", func(w http.ResponseWriter, r *http.Request) {
			technologyController.RenderCreate(w, r)
		})

		r.Post(apiVersion+"/admin/technologies", func(w http.ResponseWriter, r *http.Request) {
			technologyController.Create(w, r)
		})

		r.Get("/admin/technologies/{id}", func(w http.ResponseWriter, r *http.Request) {
			technologyController.RenderDetailEdit(w, r)
		})

		r.Put(apiVersion+"/admin/technologies/{id}", func(w http.ResponseWriter, r *http.Request) {
			technologyController.UpdateByID(w, r)
		})

		// Roles routes
		roleController := controllers.Role{
			RoleService: models.RoleService{
				DB: dbInstance,
			},
		}

		r.Get("/admin/roles", func(w http.ResponseWriter, r *http.Request) {
			roleController.RenderList(w, r)
		})

		r.Get("/admin/roles/create", func(w http.ResponseWriter, r *http.Request) {
			roleController.RenderCreate(w, r)
		})

		r.Post(apiVersion+"/admin/roles", func(w http.ResponseWriter, r *http.Request) {
			roleController.Create(w, r)
		})

		r.Get("/admin/roles/{id}", func(w http.ResponseWriter, r *http.Request) {
			roleController.RenderDetailEdit(w, r)
		})

		r.Put(apiVersion+"/admin/roles/{id}", func(w http.ResponseWriter, r *http.Request) {
			roleController.UpdateByID(w, r)
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
