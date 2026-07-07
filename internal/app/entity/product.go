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
