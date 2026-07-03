package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/RonIT-401/catalog-service/internal/app/config/section"
	"github.com/RonIT-401/catalog-service/migration"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name

	log.Printf(
		"postgres timeouts: read=%v, write=%v",
		cfg.ReadTimeout,
		cfg.WriteTimeout,
	)

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()

	sqlDB := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithDSN(dsn),
			pgdriver.WithReadTimeout(cfg.ReadTimeout),
			pgdriver.WithWriteTimeout(cfg.WriteTimeout),
		))
	sqlDB.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return &Client{
		rawBunDB: bunDB,
		cfg:      cfg,
	}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()

	if err = migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	opts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}

	m := migrate.NewMigrator(c.rawBunDB, migrations, opts...)

	if err := m.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migrations: %w", err)
	}

	applied, err := m.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	group, err := m.Migrate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to migrate: %w", err)
	}

	if len(applied) > 0 {
		oldVer, err = strconv.ParseInt(applied[0].Name, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid old migration version %q: %w", applied[0].Name, err)
		}
	}

	newVer = oldVer

	for _, mig := range group.Migrations {
		v, err := strconv.ParseInt(mig.Name, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid migration version %q: %w", mig.Name, err)
		}

		if v > newVer {
			newVer = v
		}
	}

	return oldVer, newVer, nil
}
