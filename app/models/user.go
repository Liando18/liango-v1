package models

// User represents an authenticated system user.
type User struct {
	BaseModel
	Name     string  `gorm:"type:varchar(100);not null" json:"name"`
	Email    string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string  `gorm:"type:varchar(255);not null" json:"-"`
	Role     string  `gorm:"type:varchar(50);default:'user'" json:"role"`
	Tokens   []Token `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// Token stores refresh tokens for rotating JWT auth.
type Token struct {
	BaseModel
	UserID       string `gorm:"type:uuid;not null;index" json:"user_id"`
	RefreshToken string `gorm:"type:text;not null" json:"-"`
	ExpiresAt    int64  `gorm:"not null" json:"expires_at"`
	Revoked      bool   `gorm:"default:false" json:"revoked"`
}
