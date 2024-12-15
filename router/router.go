package router

import (
	"database/sql"
	"employees-system/controllers"
	"employees-system/models"
	"employees-system/response"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

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

	r.Get(apiVersion+"/", func(w http.ResponseWriter, r *http.Request) {
		response.NewWithoutData().WithMessage("employees system api v1").Success(w)
	})

	r.Get(apiVersion+"/employees", func(w http.ResponseWriter, r *http.Request) {
		employeeController.GetAll(w, r)
	})

	r.Get(apiVersion+"/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
		employeeController.GetByID(w, r)
	})

	r.Get("/admin/employees", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/employees.html")
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
		tmpl.Execute(w, employees)
	})

	return r
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
