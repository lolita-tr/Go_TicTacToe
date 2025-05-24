package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(p DBParams) {
	db, err := sql.Open("postgres", p.URL)
	if err != nil {
		log.Fatalf("cannot open db: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("cannot create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("cannot run migration: %v", err)
	}

	if err1 := m.Up(); err1 != nil && err1 != migrate.ErrNoChange {
		log.Fatalf("migration up failed: %v", err1)
	}

	fmt.Println("Migrations applied successfully")
}
