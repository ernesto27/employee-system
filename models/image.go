package models

import "database/sql"

type Image struct {
	ID          int    `json:"id"`
	Path        string `json:"path"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type ImageService struct {
	DB          *sql.DB
	Transaction *sql.Tx
}

func (imageService *ImageService) SetTransaction(tx *sql.Tx) {
	imageService.Transaction = tx
}

func (imageService *ImageService) Create(image *Image) (int, error) {
	result, err := imageService.Transaction.Exec(`
		INSERT INTO images (path, description) VALUES (?, ?)`, image.Path, image.Description)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (imageService *ImageService) AssociateImageEmployee(imageID, employeeID int) error {
	_, err := imageService.Transaction.Exec(`
		INSERT INTO employees_images (employee_id, image_id) VALUES (?, ?)`, employeeID, imageID)
	if err != nil {
		return err
	}
	return nil
}
