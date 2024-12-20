package models

import "database/sql"

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RoleService struct {
	DB          *sql.DB
	Transaction *sql.Tx
}

func (roleService *RoleService) GetAll() ([]Role, error) {
	var roles []Role
	rows, err := roleService.DB.Query(`
		SELECT id, name FROM roles
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role Role

		err := rows.Scan(
			&role.ID, &role.Name,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (roleService *RoleService) AssociateRolEmployee(employeeID int, roleID int) error {
	_, err := roleService.Transaction.Exec(`
		INSERT INTO employees_roles (employee_id, role_id) VALUES (?, ?)
	`, employeeID, roleID)
	if err != nil {
		return err
	}

	return nil
}
