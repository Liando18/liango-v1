package repositories

import (
	"liango/app/models"
	"liango/database"

	"gorm.io/gorm"
)

// ExampleRepository handles DB operations for the Example model.
// Copy & rename this for other entities.
type ExampleRepository struct {
	db *gorm.DB
}

// NewExampleRepository creates a new instance of ExampleRepository.
func NewExampleRepository() *ExampleRepository {
	return &ExampleRepository{db: database.GetDB()}
}

// FindAll returns all examples with optional pagination.
func (r *ExampleRepository) FindAll(offset, limit int) ([]models.Example, int64, error) {
	var examples []models.Example
	var total int64

	r.db.Model(&models.Example{}).Count(&total)

	result := r.db.Offset(offset).Limit(limit).Find(&examples)
	return examples, total, result.Error
}

// FindByID returns a single example by UUID.
func (r *ExampleRepository) FindByID(id string) (*models.Example, error) {
	var example models.Example
	result := r.db.Where("id = ?", id).First(&example)
	if result.Error != nil {
		return nil, result.Error
	}
	return &example, nil
}

// Create inserts a new example into the database.
func (r *ExampleRepository) Create(example *models.Example) error {
	return r.db.Create(example).Error
}

// Update saves changes to an existing example.
func (r *ExampleRepository) Update(example *models.Example) error {
	return r.db.Save(example).Error
}

// Delete soft-deletes an example by UUID.
func (r *ExampleRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Example{}).Error
}
