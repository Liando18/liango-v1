package validations

// CreateExampleRequest is the validation struct for creating an Example.
type CreateExampleRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateExampleRequest is the validation struct for updating an Example.
type UpdateExampleRequest struct {
	Title       string `json:"title" binding:"omitempty,min=3,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// LoginRequest is the validation struct for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterRequest is the validation struct for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     string `json:"role" binding:"omitempty,oneof=admin user editor"`
}

// RefreshTokenRequest is the validation struct for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
