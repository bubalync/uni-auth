package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultMigrationsDir = "migrations"
)

func main() {
	var databaseURL, migrationsPath string

	flag.StringVar(&databaseURL, "databaseURL", "", "databaseURL")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	if databaseURL == "" {
		panic("databaseURL is required")
	}

	if migrationsPath == "" {
		migrationsPath = _defaultMigrationsDir
	}

	databaseURL += "?sslmode=disable"

	m, err := migrate.New("file://"+migrationsPath, databaseURL)
	if err != nil {
		panic(err)
	}
	err = m.Up()

	defer func() { _, _ = m.Close() }()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
