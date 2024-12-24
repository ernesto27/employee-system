package controllers

import (
	"bytes"
	"employees-system/internal/s3"
	"employees-system/models"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"
)

type ImageUploader struct {
	S3Service    s3.MyS3
	ImageService models.ImageService
}

func (iu *ImageUploader) UploadImages(r *http.Request, id int, entityType string) error {
	imagesToSave := []models.Image{}
	files := r.MultipartForm.File["files"]
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(files))

	for _, fileHeader := range files {
		wg.Add(1)
		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()

			if fileHeader.Size > (3 << 20) {
				errChan <- fmt.Errorf("file too large")
				return
			}

			file, err := fileHeader.Open()
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, file)
			if err != nil {
				errChan <- err
				return
			}

			randomString := generateRandomString(10)
			path := fmt.Sprintf("%s/%d/%s%s", entityType, id, randomString, filepath.Ext(fileHeader.Filename))
			err = iu.S3Service.Upload(&buf, path)
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			imagesToSave = append(imagesToSave, models.Image{
				Path: path,
			})
			mu.Unlock()
		}(fileHeader)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return errors.New((<-errChan).Error())
	}

	tx, err := iu.ImageService.DB.Begin()
	if err != nil {
		return errors.New(err.Error())
	}

	iu.ImageService.SetTransaction(tx)

	for _, image := range imagesToSave {
		imageID, err := iu.ImageService.Create(&image)
		if err != nil {
			tx.Rollback()
			return errors.New(err.Error())
		}

		switch entityType {
		case "employees":
			err = iu.ImageService.AssociateImageEmployee(imageID, id)
		case "projects":
			err = iu.ImageService.AssociateImageProject(imageID, id)
		default:
			err = fmt.Errorf("invalid entity type: %s", entityType)
		}

		if err != nil {
			tx.Rollback()
			return errors.New(err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.New(err.Error())
	}

	return nil
}
