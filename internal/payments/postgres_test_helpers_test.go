package payments

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newPostgresTestSchema(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration test")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("paypilot_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, `CREATE SCHEMA `+schema); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		_, _ = admin.Exec(ctx, `DROP SCHEMA `+schema+` CASCADE`)
		admin.Close()
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		_, _ = admin.Exec(ctx, `DROP SCHEMA `+schema+` CASCADE`)
		admin.Close()
		t.Fatal(err)
	}
	for _, name := range []string{"001_initial.sql", "002_reliability.sql", "003_investigations.sql"} {
		migration, err := os.ReadFile(filepath.Join("..", "..", "migrations", name))
		if err != nil {
			db.Close()
			_, _ = admin.Exec(ctx, `DROP SCHEMA `+schema+` CASCADE`)
			admin.Close()
			t.Fatal(err)
		}
		if _, err := db.Exec(ctx, string(migration)); err != nil {
			db.Close()
			_, _ = admin.Exec(ctx, `DROP SCHEMA `+schema+` CASCADE`)
			admin.Close()
			t.Fatalf("apply %s: %v", name, err)
		}
	}
	t.Cleanup(func() {
		db.Close()
		_, _ = admin.Exec(ctx, `DROP SCHEMA `+schema+` CASCADE`)
		admin.Close()
	})
	return db
}
