package router

import (
	"database/sql"
	"employees-system/controllers"
	"employees-system/models"
	"employees-system/response"
	"errors"
	"net/http"
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
