package controllers

import (
	"employees-system/models"
	"employees-system/response"
	"fmt"
	"net/http"
)

type Employee struct {
	EmployeeService models.EmployeeService
}

func (employee *Employee) GetAll(w http.ResponseWriter, r *http.Request) {
	employees, err := employee.EmployeeService.GetAll()
	if err != nil {
		fmt.Println(err)
		response.NewWithoutData().InternalServerError(w)
		return
	}

	response.New(employees).Success(w)
}
