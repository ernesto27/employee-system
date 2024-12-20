package controllers

import (
	"employees-system/models"
	"employees-system/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (employee *Employee) Create(w http.ResponseWriter, r *http.Request) {
	var newEmployee models.Employee
	jsonString := r.FormValue("jsonBody")
	err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&newEmployee)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	_, err = employee.EmployeeService.Create(newEmployee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEmployee)
}

func (employee *Employee) UpdateByID(w http.ResponseWriter, r *http.Request) {
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

	var updatedEmployee models.Employee
	jsonString := r.FormValue("jsonBody")
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&updatedEmployee)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	updatedEmployee.ID = idVal

	err = employee.EmployeeService.UpdateByID(updatedEmployee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEmployee)
}
