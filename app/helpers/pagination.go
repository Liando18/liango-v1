package helpers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"liango/app/responses"
)

type Pagination struct {
	Page    int
	PerPage int
	Offset  int
}

// GetPagination extracts and validates pagination params from query string.
// Usage: page := helpers.GetPagination(c)
func GetPagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "15"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 15
	}

	return Pagination{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

// BuildMeta builds the Meta struct for paginated responses.
func BuildMeta(p Pagination, total int64) *responses.Meta {
	totalPages := int(math.Ceil(float64(total) / float64(p.PerPage)))
	return &responses.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
