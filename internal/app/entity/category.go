package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:category"`

	ID        int64      `bun:"id,autoincrement,notnull"`
	GUID      *uuid.UUID `bun:"GUID,pk"`
	Name      string     `bun:"name,unique,notnull"`
	CreatedAt time.Time  `bun:"created_at,notnull"`
	UpdatedAt time.Time  `bun:"updated_at,notnull"`
}
