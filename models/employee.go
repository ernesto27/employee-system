package models

import (
	"database/sql"
)

type Employee struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Age           int    `json:"age"`
	Images        string `json:"images"`
	StartWorkDate string `json:"start_work_date"`
	EndWorkDate   string `json:"end_work_date"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type EmployeeService struct {
	DB *sql.DB
}

func (employeeService *EmployeeService) GetAll() ([]Employee, error) {
	var employees []Employee
	rows, err := employeeService.DB.Query(`
		SELECT 
			id, name, email, 
			age, start_work_date, end_work_date,
			created_at, updated_at 
		FROM employees`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		err := rows.Scan(
			&employee.ID, &employee.Name, &employee.Email,
			&employee.Age, &employee.StartWorkDate, &employee.EndWorkDate,
			&employee.CreatedAt, &employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}
