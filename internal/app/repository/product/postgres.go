package pproduct

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"

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

func NewRepoFromPostgres(client *rcpostgres.Client) repository.Product {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, category entity.Product) error {
	_, err := r.NewInsert().
		Model(&category).
		Exec(ctx)

	return err
}

func (r *repoPg) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error) {
	if len(guids) == 0 {
		return []entity.Product{}, entity.ErrNotFound
	}

	var products []entity.Product

	err := r.NewSelect().
		Model(&products).
		Where("GUID IN (?)", bun.List(guids)).
		Scan(ctx)

	return products, err
}

func (r *repoPg) Update(ctx context.Context, category entity.Product) error {
	result, err := r.NewUpdate().
		Model(&category).
		WherePK().
		ExcludeColumn("id", "created_at").
		Exec(ctx)

	return rcpostgres.UpdateErr(result, err)
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r.NewDelete().
		Model(&entity.Product{GUID: guid}).
		WherePK().
		Exec(ctx)

	return rcpostgres.DeleteErr(err)
}

func (r *repoPg) List(ctx context.Context, name *string, categoryGUID *uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product

	query := r.NewSelect().
		Model(&products)

	if name != nil {
		query = query.Where("name = ?", *name)
	}

	if categoryGUID != nil {
		query = query.Where("category_guid = ?", *categoryGUID)
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}
