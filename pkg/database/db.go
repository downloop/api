package database

import (
	"embed"
	"fmt"

	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrations embed.FS

const (
	host   = "downloop-downloop"
	port   = 5432
	dbname = "downloop"
)

func Init(wipe bool) (*sqlx.DB, error) {
	user := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")

	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require", user, password, host, port, dbname)
	fmt.Println(conn)

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	migration, err := migrate.NewWithInstance("iofs", source, "downloop", driver)
	if err != nil {
		return nil, err
	}

	if wipe {
		err := migration.Drop()
		if err != nil {
			return nil, err
		}
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return db, nil
}
