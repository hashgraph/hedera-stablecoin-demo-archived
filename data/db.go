package data

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"os"
)

var db *sqlx.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db = sqlx.MustConnect("pgx", os.Getenv("DATABASE_URL"))

	err = runMigrations()
	if err != nil {
		panic(err)
	}
}

func runMigrations() error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		MigrationsTable: "_migrations",
	})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://data/schema", "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		return nil
	}

	return err
}
