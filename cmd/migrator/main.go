package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "284411"
	dbname   = "yandex_smart_home"
)

func main() {
	var migrationsPath, migrationsTable, direction string
	var version uint

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.StringVar(&direction, "direction", "up", "migration direction (up/down)")
	flag.UintVar(&version, "version", 0, "specific version to migrate")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s", user, password, host, port, dbname, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	switch direction {
	case "up":
		if version > 0 {
			if err := m.Migrate(version); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					fmt.Println("no migrations to apply")

					return
				}
				panic(err)
			}
		} else {
			if err := m.Up(); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					fmt.Println("no migrations to apply")

					return
				}
				panic(err)
			}
		}
	case "down":
		if version > 0 {
			targetVersion := version - 1
			if err := m.Migrate(targetVersion); err != nil {
				panic(err)
			}
		} else {
			if err := m.Steps(-1); err != nil {
				panic(err)
			}
		}
	default:
		fmt.Println("unknown direction")
	}

	fmt.Println("migrations applied")
}
