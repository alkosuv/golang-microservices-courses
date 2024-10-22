package migration

import (
	"context"
	"github.com/pressly/goose/v3"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Up(ctx context.Context, url string, embedMigrations fs.FS) error {
	goose.SetBaseFS(embedMigrations)

	conn, err := goose.OpenDBWithDriver("pgx", url)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, conn, "migrations"); err != nil {
		return err
	}

	return nil
}
