package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminService struct {
	DB *sql.DB
}

func (adminService *AdminService) Login(email string, password string) (*Admin, error) {
	var admin Admin
	err := adminService.DB.QueryRow(`
		SELECT id, username, email, password
		FROM admins
		WHERE email = ?
	`, email).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Password,
	)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &admin, nil
}
