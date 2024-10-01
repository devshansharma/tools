package meta

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SortField struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// Pagination for limiting results and doing sorting
type Pagination struct {
	Count uint64      `json:"count"`
	Limit uint        `json:"limit"`
	Page  uint        `json:"page"`
	Sort  []SortField `json:"sort,omitempty"`
	Total uint64      `json:"total"`
}

func (p *Pagination) Parse(ctx *gin.Context) error {
	var limit = uint64(10)
	var page = uint64(1)
	var err error

	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseUint(ctx.Query("limit"), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse limit: %s", err.Error())
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.ParseUint(ctx.Query("page"), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse page: %s", err.Error())
		}
	}

	sortParam := ctx.Query("sort")
	if sortParam != "" {
		sortFields := strings.Split(sortParam, ",")
		for _, sortField := range sortFields {
			parts := strings.Split(sortField, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid sort parameter: %s", sortField)
			}

			p.Sort = append(p.Sort, SortField{
				Field: parts[0],
				Order: parts[1],
			})
		}
	}

	p.Limit = uint(limit)
	p.Page = uint(page)

	return nil
}
