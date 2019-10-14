package schema

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func MigrateUp(db *sql.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	return m.Up()
}

func MigrateDown(db *sql.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	return m.Down()
}

func newMigrator(db *sql.DB) (*migrate.Migrate, error) {
	var migrationsDir string
	if srcdir, ok := os.LookupEnv("TEST_SRCDIR"); ok {
		migrationsDir = srcdir + "/__main__/sql/migrations"
	} else if wd, ok := os.LookupEnv("BUILD_WORKING_DIRECTORY"); ok {
		migrationsDir = wd + "/sql/migrations"
	}
	log.Printf("Loading migrations from: %s", migrationsDir)
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	return migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"mysql",
		driver,
	)
}
