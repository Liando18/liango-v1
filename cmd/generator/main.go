package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// LianGo CLI Generator
// Usage:
//   go run cmd/generator/main.go make:model Product
//   go run cmd/generator/main.go make:controller Product
//   go run cmd/generator/main.go make:service Product
//   go run cmd/generator/main.go make:repository Product
//   go run cmd/generator/main.go make:crud Product

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	name := os.Args[2]

	// Capitalize first letter
	name = strings.ToUpper(name[:1]) + name[1:]

	switch command {
	case "make:model":
		generate("model", name)
	case "make:controller":
		generate("controller", name)
	case "make:service":
		generate("service", name)
	case "make:repository":
		generate("repository", name)
	case "make:crud":
		generate("model", name)
		generate("controller", name)
		generate("service", name)
		generate("repository", name)
		generate("route", name)
		generate("validation", name)
		fmt.Printf("\n✅ CRUD files generated for %s\n", name)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func generate(kind, name string) {
	tmpl, path := getTemplate(kind, name)
	if tmpl == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Printf("❌ Failed to create directory: %v\n", err)
		return
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("⚠️  File already exists: %s (skipped)\n", path)
		return
	}

	t, err := template.New(kind).Parse(tmpl)
	if err != nil {
		fmt.Printf("❌ Template error: %v\n", err)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("❌ Failed to create file: %v\n", err)
		return
	}
	defer f.Close()

	data := struct{ Name string }{Name: name}
	if err := t.Execute(f, data); err != nil {
		fmt.Printf("❌ Failed to write file: %v\n", err)
		return
	}

	fmt.Printf("✅ Created: %s\n", path)
}

func getTemplate(kind, name string) (string, string) {
	lower := strings.ToLower(name)
	_ = lower

	switch kind {
	case "model":
		return `package models

type {{.Name}} struct {
	BaseModel
	// TODO: add your fields here
	Name   string ` + "`" + `gorm:"type:varchar(255);not null" json:"name"` + "`" + `
	Status string ` + "`" + `gorm:"type:varchar(50);default:'active'" json:"status"` + "`" + `
}
`, "app/models/" + name + "Model.go"

	case "controller":
		return `package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"liango/app/responses"
	"liango/app/services"
)

type {{.Name}}Controller struct {
	service *services.{{.Name}}Service
}

func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{service: services.New{{.Name}}Service()}
}

func (ctrl *{{.Name}}Controller) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))
	items, meta, err := ctrl.service.GetAll(page, perPage)
	if err != nil {
		responses.InternalError(c, "Failed to retrieve records")
		return
	}
	responses.Paginated(c, "Records retrieved successfully", items, meta)
}

func (ctrl *{{.Name}}Controller) Show(c *gin.Context) {
	item, err := ctrl.service.GetByID(c.Param("id"))
	if err != nil {
		responses.NotFound(c, "Record not found")
		return
	}
	responses.Success(c, "Record retrieved successfully", item)
}

func (ctrl *{{.Name}}Controller) Store(c *gin.Context) {
	// TODO: bind request and call service
	responses.Created(c, "Record created successfully", nil)
}

func (ctrl *{{.Name}}Controller) Update(c *gin.Context) {
	// TODO: bind request and call service
	responses.Success(c, "Record updated successfully", nil)
}

func (ctrl *{{.Name}}Controller) Destroy(c *gin.Context) {
	if err := ctrl.service.Delete(c.Param("id")); err != nil {
		responses.NotFound(c, "Record not found")
		return
	}
	responses.Success(c, "Record deleted successfully", nil)
}
`, "app/controllers/" + name + "Controller.go"

	case "service":
		return `package services

import (
	"errors"

	"liango/app/models"
	"liango/app/repositories"
	"liango/app/responses"
	"liango/app/helpers"
)

type {{.Name}}Service struct {
	repo *repositories.{{.Name}}Repository
}

func New{{.Name}}Service() *{{.Name}}Service {
	return &{{.Name}}Service{repo: repositories.New{{.Name}}Repository()}
}

func (s *{{.Name}}Service) GetAll(page, perPage int) ([]models.{{.Name}}, *responses.Meta, error) {
	p := helpers.Pagination{Page: page, PerPage: perPage, Offset: (page - 1) * perPage}
	items, total, err := s.repo.FindAll(p.Offset, p.PerPage)
	if err != nil {
		return nil, nil, err
	}
	return items, helpers.BuildMeta(p, total), nil
}

func (s *{{.Name}}Service) GetByID(id string) (*models.{{.Name}}, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("record not found")
	}
	return item, nil
}

func (s *{{.Name}}Service) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("record not found")
	}
	return s.repo.Delete(id)
}
`, "app/services/" + name + "Service.go"

	case "repository":
		return `package repositories

import (
	"liango/app/models"
	"liango/database"

	"gorm.io/gorm"
)

type {{.Name}}Repository struct {
	db *gorm.DB
}

func New{{.Name}}Repository() *{{.Name}}Repository {
	return &{{.Name}}Repository{db: database.GetDB()}
}

func (r *{{.Name}}Repository) FindAll(offset, limit int) ([]models.{{.Name}}, int64, error) {
	var items []models.{{.Name}}
	var total int64
	r.db.Model(&models.{{.Name}}{}).Count(&total)
	result := r.db.Offset(offset).Limit(limit).Find(&items)
	return items, total, result.Error
}

func (r *{{.Name}}Repository) FindByID(id string) (*models.{{.Name}}, error) {
	var item models.{{.Name}}
	result := r.db.Where("id = ?", id).First(&item)
	return &item, result.Error
}

func (r *{{.Name}}Repository) Create(item *models.{{.Name}}) error {
	return r.db.Create(item).Error
}

func (r *{{.Name}}Repository) Update(item *models.{{.Name}}) error {
	return r.db.Save(item).Error
}

func (r *{{.Name}}Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.{{.Name}}{}).Error
}
`, "app/repositories/" + name + "Repository.go"

	case "route":
		return `package routes

import (
	"liango/app/controllers"
	"liango/app/middlewares"

	"github.com/gin-gonic/gin"
)

func Register{{.Name}}Routes(rg *gin.RouterGroup) {
	ctrl := controllers.New{{.Name}}Controller()

	group := rg.Group("/{{.Name | ToLower}}s")
	group.Use(middlewares.JWTMiddleware())
	{
		group.GET("", ctrl.Index)
		group.GET("/:id", ctrl.Show)
		group.POST("", ctrl.Store)
		group.PUT("/:id", ctrl.Update)
		group.DELETE("/:id", ctrl.Destroy)
	}
}
`, "app/routes/" + name + "Route.go"

	case "validation":
		return `package validations

type Create{{.Name}}Request struct {
	Name   string ` + "`" + `json:"name" binding:"required,min=2,max=255"` + "`" + `
	Status string ` + "`" + `json:"status" binding:"omitempty,oneof=active inactive"` + "`" + `
}

type Update{{.Name}}Request struct {
	Name   string ` + "`" + `json:"name" binding:"omitempty,min=2,max=255"` + "`" + `
	Status string ` + "`" + `json:"status" binding:"omitempty,oneof=active inactive"` + "`" + `
}
`, "app/validations/" + name + "Validation.go"
	}

	return "", ""
}

func printUsage() {
	fmt.Println(`
LianGo CLI Generator
─────────────────────────────────
Usage:
  go run cmd/generator/main.go <command> <Name>

Commands:
  make:model      <Name>   Generate a model file
  make:controller <Name>   Generate a controller file
  make:service    <Name>   Generate a service file
  make:repository <Name>   Generate a repository file
  make:crud       <Name>   Generate all CRUD files at once

Examples:
  go run cmd/generator/main.go make:model Product
  go run cmd/generator/main.go make:crud Product
`)
}
