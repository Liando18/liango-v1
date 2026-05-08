package repositories

import (
	"liango/app/models"
	"liango/database"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.GetDB()}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	result := r.db.Where("id = ?", id).First(&user)
	return &user, result.Error
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) SaveRefreshToken(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *UserRepository) FindRefreshToken(tokenStr string) (*models.Token, error) {
	var token models.Token
	result := r.db.Where("refresh_token = ? AND revoked = false", tokenStr).First(&token)
	return &token, result.Error
}

func (r *UserRepository) RevokeRefreshToken(tokenStr string) error {
	return r.db.Model(&models.Token{}).Where("refresh_token = ?", tokenStr).Update("revoked", true).Error
}
