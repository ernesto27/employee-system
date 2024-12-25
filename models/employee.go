package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type Employee struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Email         string       `json:"email"`
	DateBirth     string       `json:"dateBirth"`
	StartWorkDate string       `json:"startWorkDate"`
	EndWorkDate   string       `json:"endWorkDate"`
	Active        bool         `json:"active"`
	CreatedAt     string       `json:"created_at"`
	UpdatedAt     string       `json:"updated_at"`
	Roles         []Role       `json:"roles"`
	Technologies  []Technology `json:"technologies"`
	Projects      []Project    `json:"projects"`
	Images        []Image      `json:"images"`
}

type EmployeeService struct {
	DB *sql.DB
}

func (employeeService *EmployeeService) GetAll(page int, search string, timeRange string) ([]Employee, error) {
	offset := 0
	limit := 10

	if page > 1 {
		offset = (page - 1) * limit
	}

	query := `
		SELECT 
			id, name, email, 
			date_birth, start_work_date, end_work_date,
			created_at, updated_at 
		FROM employees`

	var args []interface{}
	var conditions []string

	if search != "" {
		searchTerm := "%" + search + "%"
		conditions = append(conditions, "(name LIKE ? OR email LIKE ?)")
		args = append(args, searchTerm, searchTerm)
	}

	if timeRange != "" && timeRange != "allTime" {
		switch timeRange {
		case "lastMonth":
			conditions = append(conditions, "start_work_date >= DATE_SUB(CURDATE(), INTERVAL 1 MONTH)")
		case "lastYear":
			conditions = append(conditions, "start_work_date >= DATE_SUB(CURDATE(), INTERVAL 1 YEAR)")
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += ` ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	var employees []Employee
	rows, err := employeeService.DB.Query(query, args...)
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
	var endWorkDate sql.NullString
	var updatedAt sql.NullString

	err := employeeService.DB.QueryRow(`
		SELECT 
			id, name, email, 
			date_birth, start_work_date, end_work_date,
			created_at, updated_at
		FROM employees WHERE id = ?`, id).Scan(
		&employee.ID, &employee.Name, &employee.Email,
		&employee.DateBirth, &employee.StartWorkDate, &endWorkDate,
		&employee.CreatedAt, &updatedAt,
	)
	if err != nil {
		return employee, err
	}

	if endWorkDate.Valid {
		employee.EndWorkDate = endWorkDate.String
	}

	if updatedAt.Valid {
		employee.UpdatedAt = updatedAt.String
	}

	roleService := RoleService{DB: employeeService.DB}
	roles, err := roleService.GetByEmployeeID(id)
	if err != nil {
		return employee, err
	}
	employee.Roles = roles

	technologyService := TechnologyService{DB: employeeService.DB}
	technologies, err := technologyService.GetByEmployeeID(id)
	if err != nil {
		return employee, err
	}
	employee.Technologies = technologies

	projectService := ProjectService{DB: employeeService.DB}
	projects, err := projectService.GetByEmployeeID(id)
	if err != nil {
		return employee, err
	}
	employee.Projects = projects

	return employee, nil
}

func (employeeService *EmployeeService) Create(employee Employee) (int, error) {
	tx, err := employeeService.DB.Begin()
	if err != nil {
		return 0, err
	}

	var endWorkDate interface{}
	if employee.EndWorkDate == "" {
		endWorkDate = nil
	} else {
		endWorkDate = employee.EndWorkDate
	}

	result, err := tx.Exec(`
		INSERT INTO employees (name, email, date_birth, start_work_date, end_work_date)
		VALUES (?, ?, ?, ?, ?)`,
		employee.Name, employee.Email, employee.DateBirth, employee.StartWorkDate, endWorkDate,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	rolService := RoleService{Transaction: tx}
	for _, role := range employee.Roles {
		err := rolService.AssociateRolEmployee(int(id), role.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	technologyService := TechnologyService{Transaction: tx}
	for _, technology := range employee.Technologies {
		err := technologyService.AssociateTechnologyEmployee(int(id), technology.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	projectService := ProjectService{Transaction: tx}
	for _, project := range employee.Projects {
		err := projectService.AssociateProjectEmployee(int(id), project.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (employeeService *EmployeeService) UpdateByID(employee Employee) error {
	tx, err := employeeService.DB.Begin()
	if err != nil {
		return err
	}

	var endWorkDate interface{}
	if employee.EndWorkDate == "" {
		endWorkDate = nil
	} else {
		endWorkDate = employee.EndWorkDate
	}

	_, err = tx.Exec(`
		UPDATE employees 
		SET name = ?, email = ?, date_birth = ?, start_work_date = ?, end_work_date = ?, active = ?
		WHERE id = ?`,
		employee.Name, employee.Email, employee.DateBirth,
		employee.StartWorkDate, endWorkDate, employee.Active,
		employee.ID,
	)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	rolService := RoleService{Transaction: tx}

	err = rolService.RemoveByEmployeeID(employee.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, role := range employee.Roles {
		err := rolService.AssociateRolEmployee(employee.ID, role.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	technologyService := TechnologyService{Transaction: tx}
	err = technologyService.RemoveByEmployeeID(employee.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, technology := range employee.Technologies {
		err := technologyService.AssociateTechnologyEmployee(employee.ID, technology.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	projectService := ProjectService{Transaction: tx}
	err = projectService.RemoveByEmployeeID(employee.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, project := range employee.Projects {
		err := projectService.AssociateProjectEmployee(employee.ID, project.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
