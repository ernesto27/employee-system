package models

import "database/sql"

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Active      bool   `json:"active"`
}

type ProjectService struct {
	DB          *sql.DB
	Transaction *sql.Tx
}

func (projectService *ProjectService) GetAll(page int) ([]Project, error) {
	var rows *sql.Rows
	var err error

	if page == -1 {
		rows, err = projectService.DB.Query(`
			SELECT id, name, description, start_date, end_date, active 
			FROM projects
			ORDER BY id DESC
		`)
	} else {
		offset := 0
		limit := 10

		if page > 1 {
			offset = (page - 1) * limit
		}

		rows, err = projectService.DB.Query(`
			SELECT id, name, description, start_date, end_date, active 
			FROM projects
			ORDER BY id DESC
			LIMIT ? OFFSET ?`,
			limit, offset,
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description,
			&project.StartDate, &project.EndDate, &project.Active,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (projectService *ProjectService) AssociateProjectEmployee(employeeID, projectID int) error {
	_, err := projectService.Transaction.Exec(`
		INSERT INTO employees_projects (employee_id, project_id)
		VALUES (?, ?)`,
		employeeID, projectID,
	)
	return err
}

func (projectService *ProjectService) GetByEmployeeID(employeeID int) ([]Project, error) {
	var projects []Project
	rows, err := projectService.DB.Query(`
		SELECT projects.id, projects.name, projects.description, projects.start_date, projects.end_date
		FROM projects
		JOIN employees_projects ON projects.id = employees_projects.project_id
		WHERE employees_projects.employee_id = ?
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var project Project

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.StartDate, &project.EndDate,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (projectService *ProjectService) RemoveByEmployeeID(employeeID int) error {
	_, err := projectService.Transaction.Exec(`
		DELETE FROM employees_projects WHERE employee_id = ?
	`, employeeID)
	if err != nil {
		return err
	}

	return nil
}
