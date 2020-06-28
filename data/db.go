package data

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"os"
)

var db *pgxpool.Pool

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db, err = pgxpool.Connect(context.TODO(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = runMigrations()
	if err != nil {
		panic(err)
	}
}

func runMigrations() error {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{
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
