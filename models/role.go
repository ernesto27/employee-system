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

func (roleService *RoleService) GetAll(page int, name string) ([]Role, error) {
	query := "SELECT id, name FROM roles"
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

	var roles []Role
	rows, err := roleService.DB.Query(query, args...)
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

func (roleService *RoleService) GetByEmployeeID(employeeID int) ([]Role, error) {
	var roles []Role
	rows, err := roleService.DB.Query(`
		SELECT roles.id, roles.name
		FROM roles
		JOIN employees_roles ON roles.id = employees_roles.role_id
		WHERE employees_roles.employee_id = ?
	`, employeeID)
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

func (roleService *RoleService) RemoveByEmployeeID(employeeID int) error {
	_, err := roleService.Transaction.Exec(`
		DELETE FROM employees_roles WHERE employee_id = ?
	`, employeeID)
	if err != nil {
		return err
	}

	return nil
}
