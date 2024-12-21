package utils

import "employees-system/models"

type TemplateData struct {
	Employees    []models.Employee
	Roles        []models.Role
	Technologies []models.Technology
	Projects     []models.Project
	Employee     models.Employee
	URL          string
}

const URL = "http://localhost:8080/admin/public"
