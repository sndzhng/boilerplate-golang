package entity

import (
	"math"
	"strings"
)

type (
	Pagination struct {
		Limit       int    `binding:"required,min=1,max=1000" form:"limit" json:"limit" gorm:"-"`
		Offset      int    `binding:"min=0" form:"offset" json:"offset" gorm:"-"`
		RecordCount *int64 `json:"record_count,omitempty" gorm:"-"`
		Total       *int   `json:"total,omitempty" gorm:"-"`
	}
	SortOrder struct {
		Sort  string `form:"sort" json:"sort,omitempty" gorm:"-"`
		Order string `form:"order" json:"order,omitempty" gorm:"-"`
	}
)

func (pagination *Pagination) CalculateTotal() {
	if pagination.RecordCount != nil {
		pagination.Total = new(int)
		*pagination.Total = int(math.Ceil(float64(*pagination.RecordCount) / float64(pagination.Limit)))
	}
}

func InitialSortOrder() SortOrder {
	return SortOrder{
		Sort:  "id",
		Order: "asc",
	}
}

func (sortOrder *SortOrder) Validate(optionalSorts ...string) bool {
	switch strings.ToLower(sortOrder.Order) {
	case "asc", "desc":
	default:
		return false
	}

	if sortOrder.Sort == "id" {
		return true
	} else {
		for _, availableSort := range optionalSorts {
			if sortOrder.Sort == availableSort {
				return true
			}
		}
		return false
	}
}
