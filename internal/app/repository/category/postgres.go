package pcategory

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/RonIT-401/catalog-service/internal/app/entity"
	"github.com/RonIT-401/catalog-service/internal/app/repository"

	rcpostgres "github.com/RonIT-401/catalog-service/internal/app/repository/conn/postgres"
)

type (
	repoPg struct {
		*_DB
	}

	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(client *rcpostgres.Client) repository.Category {
	return &repoPg{_DB: client}
}
