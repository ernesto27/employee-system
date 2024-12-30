package controllers

import (
	"bytes"
	"employees-system/internal/s3"
	"employees-system/models"
	"employees-system/response"
	"employees-system/utils"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"golang.org/x/exp/rand"
)

var funcMap = template.FuncMap{
	"hasRole": func(roleID int, employeeRoles []models.Role) bool {
		for _, r := range employeeRoles {
			if r.ID == roleID {
				return true
			}
		}
		return false
	},
	"hasTechnology": func(techID int, employeeTechs []models.Technology) bool {
		for _, t := range employeeTechs {
			if t.ID == techID {
				return true
			}
		}
		return false
	},
	"hasProject": func(projectID int, employeeProjects []models.Project) bool {
		for _, p := range employeeProjects {
			if p.ID == projectID {
				return true
			}
		}
		return false
	},
}

var tmplEmployeeList = template.Must(template.New("layout-base").Funcs(funcMap).ParseFiles("templates/layout-base.html", "templates/employees/employee-list.html"))
var tmplEmployeeCreate = template.Must(template.New("layout-base").Funcs(funcMap).ParseFiles("templates/layout-base.html", "templates/employees/employee-create.html"))
var tmplEmployeeDetailEdit = template.Must(template.New("layout-base").Funcs(funcMap).ParseFiles("templates/layout-base.html", "templates/employees/employee-detail-edit.html"))

type Employee struct {
	EmployeeService   models.EmployeeService
	RoleService       models.RoleService
	TechnologyService models.TechnologyService
	ProjectService    models.ProjectService
	S3Service         s3.MyS3
	ImageService      models.ImageService
}

func (employee *Employee) RenderList(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	searchQuery := r.URL.Query().Get("search")
	timeRange := r.URL.Query().Get("timeRange")

	employees, err := employee.EmployeeService.GetAll(page, searchQuery, timeRange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := utils.TemplateData{
		Employees: employees,
		URL:       os.Getenv("URL"),
	}

	err = tmplEmployeeList.ExecuteTemplate(w, "layout-base", data)
	if err != nil {
		fmt.Println(err)
	}
}

func (employee *Employee) GetAll(w http.ResponseWriter, r *http.Request) {
	employees, err := employee.EmployeeService.GetAll(1, "", "")
	if err != nil {
		fmt.Println(err)
		response.NewWithoutData().InternalServerError(w)
		return
	}
	response.New(employees).Success(w)
}

func (employee *Employee) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		fmt.Println("id is required")
		response.NewWithoutData().BadRequest(w)
		return
	}

	idVal, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		response.NewWithoutData().BadRequest(w)
		return
	}

	employeeData, err := employee.EmployeeService.GetByID(idVal)
	if err != nil {
		fmt.Println(err)
		response.NewWithoutData().InternalServerError(w)
		return
	}

	response.New(employeeData).Success(w)
}

func (employee *Employee) RenderCreate(w http.ResponseWriter, r *http.Request) {

	roles, err := employee.RoleService.GetAll(-1, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	technologies, err := employee.TechnologyService.GetAll(-1, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectService, err := employee.ProjectService.GetAll(-1, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := utils.TemplateData{
		URL:          os.Getenv("URL"),
		Roles:        roles,
		Technologies: technologies,
		Projects:     projectService,
	}

	err = tmplEmployeeCreate.ExecuteTemplate(w, "layout-base", data)
	if err != nil {
		fmt.Println(err)
	}
}

func (employee *Employee) Create(w http.ResponseWriter, r *http.Request) {
	var newEmployee models.Employee
	jsonString := r.FormValue("jsonBody")
	err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&newEmployee)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	id, err := employee.EmployeeService.Create(newEmployee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = uploadImages(r, id, employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEmployee)
}

func (employee *Employee) RenderDetailEdit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	employeeByID, err := employee.EmployeeService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roles, err := employee.RoleService.GetAll(-1, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	technologies, err := employee.TechnologyService.GetAll(-1, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projects, err := employee.ProjectService.GetAll(-1, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	images, err := employee.ImageService.GetImagesByEmployeeID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employeeByID.Images = images

	data := utils.TemplateData{
		Employee:     employeeByID,
		Roles:        roles,
		Technologies: technologies,
		Projects:     projects,
		URL:          os.Getenv("URL"),
	}

	err = tmplEmployeeDetailEdit.ExecuteTemplate(w, "layout-base", data)
	if err != nil {
		fmt.Println(err)
	}
}

func (employee *Employee) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		fmt.Println("id is required")
		response.NewWithoutData().BadRequest(w)
		return
	}

	idVal, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		response.NewWithoutData().BadRequest(w)
		return
	}

	var updatedEmployee models.Employee
	jsonString := r.FormValue("jsonBody")
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&updatedEmployee)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	updatedEmployee.ID = idVal
	err = employee.EmployeeService.UpdateByID(updatedEmployee)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = uploadImages(r, idVal, employee)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("imagesToDelete") != "" {
		imagesIDsToDelete := strings.Split(r.FormValue("imagesToDelete"), ",")
		for _, imageID := range imagesIDsToDelete {
			imageIDVal, err := strconv.Atoi(imageID)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = employee.ImageService.DeleteByID(imageIDVal, "employees_images")
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEmployee)
}

func uploadImages(r *http.Request, id int, employee *Employee) error {
	imagesToSave := []models.Image{}
	files := r.MultipartForm.File["files"]
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(files))

	for _, fileHeader := range files {
		wg.Add(1)
		randomString := generateRandomString(10)
		go func(fileHeader *multipart.FileHeader, randomString string) {
			defer wg.Done()

			if fileHeader.Size > (3 << 20) {
				errChan <- fmt.Errorf("file too large")
				return
			}

			fmt.Println("Uploading file", fileHeader.Filename)

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

			path := fmt.Sprintf("employees/%d/%s.%s", id, randomString, filepath.Ext(fileHeader.Filename))
			err = employee.S3Service.Upload(&buf, path)
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			imagesToSave = append(imagesToSave, models.Image{
				Path: path,
			})
			mu.Unlock()
		}(fileHeader, randomString)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return errors.New((<-errChan).Error())
	}

	tx, err := employee.ImageService.DB.Begin()
	if err != nil {
		return errors.New(err.Error())
	}

	employee.ImageService.SetTransaction(tx)

	for _, image := range imagesToSave {
		imageID, err := employee.ImageService.Create(&image)
		if err != nil {
			tx.Rollback()
			return errors.New(err.Error())
		}

		err = employee.ImageService.AssociateImageEmployee(imageID, id)
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

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
