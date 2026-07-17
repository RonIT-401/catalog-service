package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:category"`

	ID        int64     `bun:"id,autoincrement,notnull"`
	GUID      uuid.UUID `bun:"guid,pk"`
	Name      string    `bun:"name,unique,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type RequestCategoryCreate struct {
	Name string `json:"name" binding:"required,min=2,max=255"`
}

type RequestCategoryUpdate struct {
	Name string `json:"name" binding:"required,min=2,max=255"`
}

type ResponseCategoryCreate struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseCategoryUpdate struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseCategoryList struct {
	Data []ResponseCategoryListItem `json:"data"`
}

type ResponseCategoryListItem struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
