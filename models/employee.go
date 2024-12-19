package models

import (
	"database/sql"
)

type Employee struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	DateBirth     string `json:"dateBirth"`
	Image         string `json:"images"`
	StartWorkDate string `json:"startWorkDate"`
	EndWorkDate   string `json:"endWorkDate"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type EmployeeService struct {
	DB *sql.DB
}

func (employeeService *EmployeeService) GetAll(page int) ([]Employee, error) {
	offset := 0
	limit := 10

	if page > 1 {
		offset = (page - 1) * limit
	}

	var employees []Employee
	rows, err := employeeService.DB.Query(`
		SELECT 
			id, name, email, 
			date_birth, start_work_date, end_work_date,
			created_at, updated_at 
		FROM employees
		ORDER BY id DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		var endWorkDate sql.NullString
		var updatedAt sql.NullString

		err := rows.Scan(
			&employee.ID, &employee.Name, &employee.Email,
			&employee.DateBirth, &employee.StartWorkDate, &endWorkDate,
			&employee.CreatedAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endWorkDate.Valid {
			employee.EndWorkDate = endWorkDate.String
		}

		if updatedAt.Valid {
			employee.UpdatedAt = updatedAt.String
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

func (employeeService *EmployeeService) GetByID(id int) (Employee, error) {
	var employee Employee
	err := employeeService.DB.QueryRow(`
		SELECT 
			id, name, email, 
			age, start_work_date, end_work_date,
			created_at, updated_at, image
		FROM employees WHERE id = ?`, id).Scan(
		&employee.ID, &employee.Name, &employee.Email,
		&employee.DateBirth, &employee.StartWorkDate, &employee.EndWorkDate,
		&employee.CreatedAt, &employee.UpdatedAt, &employee.Image,
	)
	if err != nil {
		return employee, err
	}

	return employee, nil
}

func (employeeService *EmployeeService) Create(employee Employee) (int, error) {
	result, err := employeeService.DB.Exec(`
		INSERT INTO employees (name, email, date_birth, start_work_date, image)
		VALUES (?, ?, ?, ?, ?)`,
		employee.Name, employee.Email, employee.DateBirth, employee.StartWorkDate, employee.Image,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
