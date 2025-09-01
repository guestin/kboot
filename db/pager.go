package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/guestin/mob/merrors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PageableTable interface {
	schema.Tabler
	OrderColLimit() map[string]string
	Filter() []string
}

type PageRequest struct {
	Page     *int `json:"page" query:"page" form:"page" path:"page" form:"page" validate:"omitempty,gt=0"`
	PageSize *int `json:"pageSize" query:"pageSize" form:"pageSize" path:"pageSize" form:"pageSize" validate:"omitempty,gt=0"`

	Begin int64  `json:"begin" query:"begin" form:"begin" path:"begin" form:"begin" validate:"gte=0"`
	End   *int64 `json:"end" query:"end" form:"end" path:"end" form:"end" validate:"omitempty,gtfield=Begin"`

	Key string `json:"key" query:"key" form:"key" path:"key" form:"key"`

	OrderBy string `json:"orderBy" query:"orderBy" form:"orderBy" path:"orderBy" form:"orderBy"`

	Order string `json:"order" query:"order" form:"order" path:"order" form:"order" validate:"omitempty,oneof=ASC DESC"`
}

func (this *PageRequest) PageV() int {
	if this.Page != nil && *this.Page > 0 {
		return *this.Page
	}
	return 1
}

func (this *PageRequest) PageSizeV() int {
	if this.PageSize != nil && *this.PageSize > 0 {
		return *this.PageSize
	}
	return 10
}

func (this *PageRequest) BeginV() int64 {
	if this.Begin > 0 {
		return this.Begin
	}
	return 0
}
func (this *PageRequest) EndV() int64 {
	if this.End != nil && *this.End > 0 {
		return *this.End
	}
	return 0
}

func (this *PageRequest) OrderV() string {
	if this.Order != "" {
		return this.Order
	}
	return "ASC"
}

func (this *PageRequest) Offset() int {
	return (this.PageV() - 1) * this.PageSizeV()
}

func (this *PageRequest) Limit() int {
	return this.PageSizeV()
}

func (this *PageRequest) BuildResponse(results interface{}) *PageResponse {
	return &PageResponse{
		Total:    0,
		Page:     this.PageV(),
		PageSize: this.PageSizeV(),
		Results:  results,
	}
}

type PageResponse struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Results  interface{} `json:"results"`
}

//goland:noinspection ALL
func PageQuery(query *gorm.DB, page PageRequest, table PageableTable, defaultOrder string, optFilterCol ...string) (*gorm.DB, error) {
	query = query.Where(fmt.Sprintf("%s.created_at >= ?", table.TableName()), time.Unix(page.BeginV(), 0))
	if page.EndV() > 0 {
		query = query.Where(fmt.Sprintf("%s.created_at <= ?", table.TableName()), time.Unix(page.EndV(), 0))
	}
	if len(page.OrderBy) > 0 {
		//check order
		colLimit := table.OrderColLimit()
		orderCol, ok := colLimit[page.OrderBy]
		if !ok {
			return query, merrors.Errorf("orderBy '%s' not allowed , must be one of [%s]", page.OrderBy,
				mkArrayString(colLimit))
		}
		query = query.Order(fmt.Sprintf("%s %s", orderCol, page.OrderV()))
	} else if len(defaultOrder) > 0 {
		query = query.Order(defaultOrder[0])
	} else {
		//default order by created desc
		query = query.Order(fmt.Sprintf("%s.created_at DESC", table.TableName()))
	}
	if len(page.Key) > 0 {
		key := page.Key
		orCols := make([]string, 0)
		args := make([]interface{}, 0)
		for _, col := range table.Filter() {
			orCols = append(orCols, fmt.Sprintf("%s LIKE ? ", col))
			args = append(args, "%"+key+"%")
		}
		for _, col := range optFilterCol {
			orCols = append(orCols, fmt.Sprintf("%s LIKE ? ", col))
			args = append(args, "%"+key+"%")
		}
		if len(orCols) > 0 {
			orQueryStr := fmt.Sprintf("(%s)", strings.Join(orCols, " OR "))
			query = query.Where(orQueryStr, args...)
		}
	}
	return query, nil
}

func mkArrayString(names map[string]string) string {
	dbStrElements := make([]string, 0, len(names))
	for k := range names {
		dbStrElements = append(dbStrElements, fmt.Sprintf("'%s'", k))
	}
	return strings.Join(dbStrElements, ",")
}
