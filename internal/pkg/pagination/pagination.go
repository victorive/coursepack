package pagination

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

type Query struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type PaginatedResponse[T any] struct {
	Items      []T `json:"items"`
	Total      int64 `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalPages int64 `json:"totalPages"`
}

func Parse(c *fiber.Ctx) *Query {
	pageStr := c.Query("page", strconv.Itoa(DefaultPage))
	perPageStr := c.Query("perPage", strconv.Itoa(DefaultPerPage))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = DefaultPage
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return &Query{
		Page:    page,
		PerPage: perPage,
	}
}

func (q *Query) Offset() int {
	if q.Page < 1 {
		q.Page = 1
	}
	return (q.Page - 1) * q.PerPage
}

func (q *Query) Limit() int {
	if q.PerPage < 1 {
		q.PerPage = DefaultPerPage
	}
	if q.PerPage > MaxPerPage {
		q.PerPage = MaxPerPage
	}
	return q.PerPage
}

func NewPaginated[T any](items []T, total int64, q *Query) *PaginatedResponse[T] {
	if items == nil {
		items = []T{}
	}

	totalPages := (total + int64(q.PerPage) - 1) / int64(q.PerPage)
	if totalPages < 1 {
		totalPages = 1
	}

	return &PaginatedResponse[T]{
		Items:      items,
		Total:      total,
		Page:       q.Page,
		PerPage:    q.PerPage,
		TotalPages: totalPages,
	}
}
