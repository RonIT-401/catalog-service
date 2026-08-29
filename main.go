package main

import (
	"context"
	"os"

	"github.com/RonIT-401/catalog-service/internal/app/config"
	hcategory "github.com/RonIT-401/catalog-service/internal/app/handler/http/category"
	rhealth "github.com/RonIT-401/catalog-service/internal/app/handler/http/health"
	hproduct "github.com/RonIT-401/catalog-service/internal/app/handler/http/product"
	rprocessor "github.com/RonIT-401/catalog-service/internal/app/processor/http"
	pcategory "github.com/RonIT-401/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/RonIT-401/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/RonIT-401/catalog-service/internal/app/repository/product"
	scategory "github.com/RonIT-401/catalog-service/internal/app/service/category"
	sproduct "github.com/RonIT-401/catalog-service/internal/app/service/product"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	config.Load(config.LoadArgs{
		Output:          os.Stdout,
		EnableSimpleLog: true,
	})

	cfg := config.Root

	pgClient, err := rcpostgres.NewClient(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}

	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
	if oldVer != newVer {
		log.Info().Int64("old_version", oldVer).Int64("new_version", newVer).Msg("Database migrated")
	} else {
		log.Info().Int64("version", newVer).Msg("Database is up to date")
	}

	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	productRepo := pproduct.NewRepoFromPostgres(pgClient)

	categorySvc := scategory.NewService(categoryRepo, productRepo)
	productSvc := sproduct.NewService(productRepo, categoryRepo)

	hHealth := rhealth.NewHandler()
	hCategory := hcategory.NewHandler(categorySvc)
	hProduct := hproduct.NewHandler(productSvc)

	httpServer := rprocessor.NewHTTP(hHealth, hCategory, hProduct, cfg.Processor.WebServer)
	if err := httpServer.Serve(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
