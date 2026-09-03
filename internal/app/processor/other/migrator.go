package pprocessor

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/RonIT-401/catalog-service/internal/app/processor"
	"github.com/RonIT-401/catalog-service/internal/app/repository"
)

type procMigrate struct {
	migrator repository.Migrate
}

func NewMigrator(migrator repository.Migrate) processor.Processor {
	return &procMigrate{
		migrator: migrator,
	}
}

func (p *procMigrate) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	processor.Wrap(ctx, wg, p.job)
}

func (p *procMigrate) job(ctx context.Context) {
	oldVer, newVer, err := p.migrator.Migrate((ctx))
	if err != nil {
		log.Error().Err(err).Msg("failed to run database migration")
		return
	}

	if oldVer != newVer {
		log.Info().Int64("old version", oldVer).
			Int64("new version", newVer).
			Msg("schema has been updated")
	}
	if oldVer == newVer {
		log.Info().Int64("version", oldVer).Msg("schema is up to date")
	}
}
