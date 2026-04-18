package pagination

import (
	"fmt"
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
	PerPage int `json:"per_page"`
}

type Links struct {
	First *string `json:"first"`
	Last  *string `json:"last"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

type Meta struct {
	Total       int64 `json:"total"`
	Count       int   `json:"count"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
	Links       Links `json:"links"`
}

type PaginatedResponse[T any] struct {
	Items []T  `json:"items"`
	Meta  Meta `json:"meta"`
}

func Parse(c *fiber.Ctx) *Query {
	pageStr := c.Query("page", strconv.Itoa(DefaultPage))
	perPageStr := c.Query("per_page", strconv.Itoa(DefaultPerPage))

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

func NewPaginated[T any](c *fiber.Ctx, items []T, total int64, q *Query) *PaginatedResponse[T] {
	if items == nil {
		items = []T{}
	}

	totalPages := (total + int64(q.PerPage) - 1) / int64(q.PerPage)
	if totalPages < 1 {
		totalPages = 1
	}

	first := pageURL(c, 1, q.PerPage)
	last := pageURL(c, int(totalPages), q.PerPage)

	var prev *string
	if q.Page > 1 {
		u := pageURL(c, q.Page-1, q.PerPage)
		prev = &u
	}

	var next *string
	if int64(q.Page) < totalPages {
		u := pageURL(c, q.Page+1, q.PerPage)
		next = &u
	}

	return &PaginatedResponse[T]{
		Items: items,
		Meta: Meta{
			Total:       total,
			Count:       len(items),
			PerPage:     q.PerPage,
			CurrentPage: q.Page,
			TotalPages:  totalPages,
			Links: Links{
				First: &first,
				Last:  &last,
				Prev:  prev,
				Next:  next,
			},
		},
	}
}

func pageURL(c *fiber.Ctx, page, perPage int) string {
	return fmt.Sprintf("%s%s?page=%d&per_page=%d", c.BaseURL(), c.Path(), page, perPage)
}
