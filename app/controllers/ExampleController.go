package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"liango/app/responses"
	"liango/app/services"
	"liango/app/validations"
)

// ExampleController handles HTTP requests for the Example resource.
// Copy & rename this for other entities.
type ExampleController struct {
	service *services.ExampleService
}

// NewExampleController creates a new ExampleController.
func NewExampleController() *ExampleController {
	return &ExampleController{service: services.NewExampleService()}
}

// Index returns a paginated list of examples.
// GET /examples
func (ctrl *ExampleController) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	examples, meta, err := ctrl.service.GetAll(page, perPage)
	if err != nil {
		responses.InternalError(c, "Failed to retrieve examples")
		return
	}

	responses.Paginated(c, "Examples retrieved successfully", examples, meta)
}

// Show returns a single example by ID.
// GET /examples/:id
func (ctrl *ExampleController) Show(c *gin.Context) {
	id := c.Param("id")

	example, err := ctrl.service.GetByID(id)
	if err != nil {
		responses.NotFound(c, "Example not found")
		return
	}

	responses.Success(c, "Example retrieved successfully", example)
}

// Store creates a new example.
// POST /examples
func (ctrl *ExampleController) Store(c *gin.Context) {
	var req validations.CreateExampleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.UnprocessableEntity(c, "Validation failed", err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	example, err := ctrl.service.Create(req, userID.(string))
	if err != nil {
		responses.InternalError(c, "Failed to create example")
		return
	}

	responses.Created(c, "Example created successfully", example)
}

// Update modifies an existing example.
// PUT /examples/:id
func (ctrl *ExampleController) Update(c *gin.Context) {
	id := c.Param("id")

	var req validations.UpdateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.UnprocessableEntity(c, "Validation failed", err.Error())
		return
	}

	example, err := ctrl.service.Update(id, req)
	if err != nil {
		if err.Error() == "example not found" {
			responses.NotFound(c, "Example not found")
			return
		}
		responses.InternalError(c, "Failed to update example")
		return
	}

	responses.Success(c, "Example updated successfully", example)
}

// Destroy soft-deletes an example.
// DELETE /examples/:id
func (ctrl *ExampleController) Destroy(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.service.Delete(id); err != nil {
		if err.Error() == "example not found" {
			responses.NotFound(c, "Example not found")
			return
		}
		responses.InternalError(c, "Failed to delete example")
		return
	}

	responses.Success(c, "Example deleted successfully", nil)
}
