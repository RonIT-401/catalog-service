package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product"`

	ID           int64     `bun:"id,autoincrement,notnull"`
	GUID         uuid.UUID `bun:"guid,pk"`
	Name         string    `bun:"name,unique,notnull"`
	Description  *string   `bun:"description"`
	Price        int64     `bun:"price,notnull"`
	CategoryGUID uuid.UUID `bun:"category_guid,notnull"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt    time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type RequestProductCreate struct {
	Name         string    `json:"name" binding:"required,min=2,max=255"`
	Description  *string   `json:"description" binding:"omitempty,max=1000"`
	Price        int64     `json:"price" binding:"required,gt=0"`
	CategoryGUID uuid.UUID `json:"category_guid" binding:"required,uuid"`
}

type RequestProductUpdate struct {
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        *int64    `json:"price" binding:"required,gt=0"`
	CategoryGUID uuid.UUID `json:"category_guid"`
}

type RequestProductList struct {
	CategoryGUID *uuid.UUID `json:"category_guid"`
}

type ResponseProductCreate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        int64     `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
}

type ResponseProductUpdate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        int64     `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ResponseProductList struct {
	Data []ResponseProductListItem `json:"data"`
}

type ResponseProductListItem struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Price        int64     `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
