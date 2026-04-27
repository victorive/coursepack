package pagination

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const DefaultCursorLimit = 20

type CursorQuery struct {
	After string
	Limit int
}

type CursorMeta struct {
	NextCursor *string `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
	Count      int     `json:"count"`
}

type CursorResponse[T any] struct {
	Items []T        `json:"items"`
	Meta  CursorMeta `json:"meta"`
}

func ParseCursor(c *fiber.Ctx) *CursorQuery {
	after := c.Query("after", "")

	limit := DefaultCursorLimit
	if s := c.Query("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > MaxPerPage {
		limit = MaxPerPage
	}

	return &CursorQuery{After: after, Limit: limit}
}

func NewCursorResponse[T any](items []T, nextCursor *string) *CursorResponse[T] {
	if items == nil {
		items = []T{}
	}
	return &CursorResponse[T]{
		Items: items,
		Meta: CursorMeta{
			NextCursor: nextCursor,
			HasMore:    nextCursor != nil,
			Count:      len(items),
		},
	}
}
