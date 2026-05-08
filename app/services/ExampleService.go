package services

import (
	"errors"

	"liango/app/helpers"
	"liango/app/models"
	"liango/app/repositories"
	"liango/app/responses"
	"liango/app/validations"
)

// ExampleService handles business logic for the Example entity.
// Copy & rename this for other entities.
type ExampleService struct {
	repo *repositories.ExampleRepository
}

// NewExampleService creates a new ExampleService.
func NewExampleService() *ExampleService {
	return &ExampleService{
		repo: repositories.NewExampleRepository(),
	}
}

// GetAll returns paginated examples.
func (s *ExampleService) GetAll(page, perPage int) ([]models.Example, *responses.Meta, error) {
	p := helpers.Pagination{Page: page, PerPage: perPage, Offset: (page - 1) * perPage}
	examples, total, err := s.repo.FindAll(p.Offset, p.PerPage)
	if err != nil {
		return nil, nil, err
	}
	meta := helpers.BuildMeta(p, total)
	return examples, meta, nil
}

// GetByID returns a single example or error if not found.
func (s *ExampleService) GetByID(id string) (*models.Example, error) {
	example, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("example not found")
	}
	return example, nil
}

// Create validates and creates a new example.
func (s *ExampleService) Create(req validations.CreateExampleRequest, userID string) (*models.Example, error) {
	example := &models.Example{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		UserID:      userID,
	}

	if example.Status == "" {
		example.Status = "active"
	}

	if err := s.repo.Create(example); err != nil {
		return nil, err
	}

	return example, nil
}

// Update applies partial updates to an existing example.
func (s *ExampleService) Update(id string, req validations.UpdateExampleRequest) (*models.Example, error) {
	example, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("example not found")
	}

	if req.Title != "" {
		example.Title = req.Title
	}
	if req.Description != "" {
		example.Description = req.Description
	}
	if req.Status != "" {
		example.Status = req.Status
	}

	if err := s.repo.Update(example); err != nil {
		return nil, err
	}

	return example, nil
}

// Delete soft-deletes an example.
func (s *ExampleService) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("example not found")
	}
	return s.repo.Delete(id)
}
