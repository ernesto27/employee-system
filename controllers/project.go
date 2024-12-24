package controllers

import (
	"employees-system/internal/s3"
	"employees-system/models"
	"employees-system/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type Project struct {
	ProjectService models.ProjectService
	S3Service      s3.MyS3
	ImageService   models.ImageService
}

var tmplList = template.Must(template.ParseFiles("templates/layout-base.html", "templates/projects/project-list.html"))

func (project *Project) RenderList(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	nameQuery := r.URL.Query().Get("name")

	projects, err := project.ProjectService.GetAll(page, nameQuery)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := utils.TemplateData{
		URL:      utils.URL,
		Projects: projects,
	}

	err = tmplList.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}

var tmplCreate = template.Must(
	template.ParseFiles("templates/layout-base.html", "templates/projects/project-create.html", "templates/components/upload-file.html"),
)

func (project *Project) RenderCreate(w http.ResponseWriter, r *http.Request) {
	templateData := utils.TemplateData{
		URL: utils.URL,
	}

	err := tmplCreate.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}

var tmplDetailEdit = template.Must(template.ParseFiles("templates/layout-base.html", "templates/projects/project-detail-edit.html"))

func (project *Project) RenderDetailEdit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	projectDetail, err := project.ProjectService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	images, err := project.ImageService.GetImagesByProjectID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	projectDetail.Images = images

	templateData := utils.TemplateData{
		URL:     utils.URL,
		Project: projectDetail,
	}

	err = tmplDetailEdit.ExecuteTemplate(w, "layout-base", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}

func (project *Project) UpdateByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var newProject models.Project
	jsonString := r.FormValue("jsonBody")
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&newProject)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	newProject.ID = id

	err = project.ProjectService.UpdateByID(newProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploader := &ImageUploader{
		S3Service:    project.S3Service,
		ImageService: project.ImageService,
	}

	err = uploader.UploadImages(r, id, "projects")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imagesIDsToDelete := strings.Split(r.FormValue("imagesToDelete"), ",")
	for _, imageID := range imagesIDsToDelete {
		imageIDVal, err := strconv.Atoi(imageID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = project.ImageService.DeleteByID(imageIDVal, "projects_images")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProject)
}

func (project *Project) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	var newProject models.Project
	jsonString := r.FormValue("jsonBody")
	err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&newProject)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	insertID, err := project.ProjectService.Create(newProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploader := &ImageUploader{
		S3Service:    project.S3Service,
		ImageService: project.ImageService,
	}

	err = uploader.UploadImages(r, insertID, "projects")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProject)
}
