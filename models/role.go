package models

import "database/sql"

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

func (roleService *RoleService) GetByID(id int) (Role, error) {
	var role Role
	var description sql.NullString
	err := roleService.DB.QueryRow("SELECT id, name, description FROM roles WHERE id = ?", id).
		Scan(&role.ID, &role.Name, &description)
	if err != nil {
		return role, err
	}

	if description.Valid {
		role.Description = description.String
	}
	return role, nil
}

func (roleService *RoleService) Create(role Role) (int, error) {
	res, err := roleService.DB.Exec("INSERT INTO roles (name, description) VALUES (?, ?)", role.Name, role.Description)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (roleService *RoleService) UpdateByID(id int, role Role) error {
	_, err := roleService.DB.Exec(`
		UPDATE roles 
		SET name = ?, description = ?
		WHERE id = ?
	`, role.Name, role.Description, id)
	if err != nil {
		return err
	}

	return nil
}
