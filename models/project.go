package models

import "database/sql"

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Active      bool   `json:"active"`
	Links       string `json:"links"`
	Contacts    string `json:"contacts"`
}

type ProjectService struct {
	DB          *sql.DB
	Transaction *sql.Tx
}

func (projectService *ProjectService) GetAll(page int, searchName string) ([]Project, error) {
	var rows *sql.Rows
	var err error
	var query string
	var args []interface{}

	baseQuery := `
		SELECT id, name, description, start_date, end_date, active 
		FROM projects
		WHERE 1=1`

	if searchName != "" {
		baseQuery += ` AND name LIKE ?`
		args = append(args, "%"+searchName+"%")
	}

	baseQuery += ` ORDER BY id DESC`

	if page == -1 {
		query = baseQuery
		rows, err = projectService.DB.Query(query, args...)
	} else {
		offset := 0
		limit := 10

		if page > 1 {
			offset = (page - 1) * limit
		}

		query = baseQuery + ` LIMIT ? OFFSET ?`
		args = append(args, limit, offset)
		rows, err = projectService.DB.Query(query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		var endDate sql.NullString

		err := rows.Scan(
			&project.ID, &project.Name, &project.Description,
			&project.StartDate, &endDate, &project.Active,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			project.EndDate = endDate.String
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (projectService *ProjectService) Create(project Project) error {
	var endDate interface{}
	if project.EndDate == "" {
		endDate = nil
	} else {
		endDate = project.EndDate
	}

	_, err := projectService.DB.Exec(`
		INSERT INTO projects (name, description, start_date, end_date, active, links, contacts)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		project.Name, project.Description,
		project.StartDate, endDate, project.Active,
		project.Links, project.Contacts,
	)
	if err != nil {
		return err
	}

	return nil
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
