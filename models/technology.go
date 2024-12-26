package models

import "database/sql"

type Technology struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TechnologyService struct {
	DB          *sql.DB
	Transaction *sql.Tx
}

func (technologyService *TechnologyService) GetAll(page int, name string) ([]Technology, error) {
	query := "SELECT id, name FROM technologies"
	var args []interface{}

	if name != "" {
		query += " WHERE name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if page != -1 {
		offset := 0
		limit := 10
		if page > 1 {
			offset = (page - 1) * limit
		}
		query += " ORDER BY id DESC LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	var technologies []Technology
	rows, err := technologyService.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var technology Technology

		err := rows.Scan(
			&technology.ID, &technology.Name,
		)
		if err != nil {
			return nil, err
		}

		technologies = append(technologies, technology)
	}

	return technologies, nil
}

func (technologyService *TechnologyService) AssociateTechnologyEmployee(employeeID int, technologyID int) error {
	_, err := technologyService.Transaction.Exec(`
		INSERT INTO employees_technologies (employee_id, technology_id) VALUES (?, ?)
	`, employeeID, technologyID)
	if err != nil {
		return err
	}

	return nil
}

func (technologyService *TechnologyService) GetByEmployeeID(employeeID int) ([]Technology, error) {
	var technologies []Technology
	rows, err := technologyService.DB.Query(`
		SELECT technologies.id, technologies.name
		FROM technologies
		JOIN employees_technologies ON technologies.id = employees_technologies.technology_id
		WHERE employees_technologies.employee_id = ?
	`, employeeID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var technology Technology

		err := rows.Scan(
			&technology.ID, &technology.Name,
		)
		if err != nil {
			return nil, err
		}

		technologies = append(technologies, technology)
	}

	return technologies, nil
}

func (technologyService *TechnologyService) RemoveByEmployeeID(employeeID int) error {
	_, err := technologyService.Transaction.Exec(`
		DELETE FROM employees_technologies WHERE employee_id = ?
	`, employeeID)
	if err != nil {
		return err
	}

	return nil
}

func (technologyService *TechnologyService) Create(technology Technology) (int, error) {
	result, err := technologyService.DB.Exec(`
		INSERT INTO technologies (name) VALUES (?)
	`, technology.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (technologyService *TechnologyService) UpdateByID(id int, technology Technology) error {
	_, err := technologyService.DB.Exec(`
		UPDATE technologies 
		SET name = ?, description = ?
		WHERE id = ?
	`, technology.Name, technology.Description, id)

	return err
}

func (technologyService *TechnologyService) GetByID(id int) (Technology, error) {
	var technology Technology
	var description sql.NullString

	err := technologyService.DB.QueryRow(`
		SELECT id, name, description
		FROM technologies 
		WHERE id = ?
	`, id).Scan(&technology.ID, &technology.Name, &description)

	if description.Valid {
		technology.Description = description.String
	}

	return technology, err
}
