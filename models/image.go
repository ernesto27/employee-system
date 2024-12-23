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

func (ImageService *ImageService) GetImagesByEmployeeID(employeeID int) ([]Image, error) {
	rows, err := ImageService.DB.Query(`
		SELECT images.id, images.path, images.description, images.created_at
		FROM images
		INNER JOIN employees_images ON images.id = employees_images.image_id
		WHERE employees_images.employee_id = ?`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []Image{}
	for rows.Next() {
		image := Image{}
		err := rows.Scan(&image.ID, &image.Path, &image.Description, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func (imageService *ImageService) DeleteByID(id int) error {
	tx, err := imageService.DB.Begin()
	if err != nil {
		return err
	}

	imageService.SetTransaction(tx)

	_, err = imageService.Transaction.Exec(`DELETE FROM images WHERE id = ?`, id)
	if err != nil {
		imageService.Transaction.Rollback()
		return err
	}

	_, err = imageService.Transaction.Exec(`DELETE FROM employees_images WHERE image_id = ?`, id)
	if err != nil {
		imageService.Transaction.Rollback()
		return err
	}

	err = imageService.Transaction.Commit()
	if err != nil {
		imageService.Transaction.Rollback()
		return err
	}

	return nil
}
