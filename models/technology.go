package models

import "database/sql"

type Technology struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TechnologyService struct {
	DB          *sql.DB
	Transaction *sql.Tx
}

func (technologyService *TechnologyService) GetAll() ([]Technology, error) {
	var technologies []Technology
	rows, err := technologyService.DB.Query(`
		SELECT id, name FROM technologies
	`)
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