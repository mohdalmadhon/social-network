package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() error {
	path := os.Getenv("ORBIT_DB_PATH")
	if path == "" {
		path = "db/social_network.db"
	}
	m, err := migrate.New(
		"file://internal/migrations",
		"sqlite3://"+path,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	fmt.Println("Migrations applied successfully")
	return nil
}
