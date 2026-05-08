package models

// Example is a template model — copy & rename for your own entity.
// It embeds BaseModel for UUID PK + timestamps + soft delete.
type Example struct {
	BaseModel
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"type:varchar(50);default:'active'" json:"status"`
	UserID      string `gorm:"type:uuid;index" json:"user_id"`
}
