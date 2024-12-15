package controllers

import (
	"employees-system/models"
	"employees-system/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Employee struct {
	EmployeeService models.EmployeeService
}

func (employee *Employee) GetAll(w http.ResponseWriter, r *http.Request) {
	employees, err := employee.EmployeeService.GetAll(1)
	if err != nil {
		fmt.Println(err)
		response.NewWithoutData().InternalServerError(w)
		return
	}
	response.New(employees).Success(w)
}

func (employee *Employee) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		fmt.Println("id is required")
		response.NewWithoutData().BadRequest(w)
		return
	}

	idVal, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		response.NewWithoutData().BadRequest(w)
		return
	}

	employeeData, err := employee.EmployeeService.GetByID(idVal)
	if err != nil {
		fmt.Println(err)
		response.NewWithoutData().InternalServerError(w)
		return
	}

	response.New(employeeData).Success(w)

}
