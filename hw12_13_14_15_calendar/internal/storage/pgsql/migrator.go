package pgsqlstorage

import (
	"database/sql"
	"fmt"

	_ "github.com/ekhvalov/otus-golang/hw12_13_14_15_calendar/migrations/pgsql" // init migrations
	"github.com/hashicorp/go-multierror"
	_ "github.com/jackc/pgx/v5/stdlib" // init pgx stdlib driver
	"github.com/pressly/goose/v3"
)

func NewMigrator(conf Config) Migrator {
	return Migrator{conf: conf}
}

type Migrator struct {
	conf Config
}

func (m Migrator) Run(command string) (err error) {
	db, err := sql.Open("pgx", m.conf.GetDSN())
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer func() {
		dbErr := db.Close()
		if dbErr != nil {
			err = multierror.Append(err, dbErr)
		}
	}()

	err = goose.Run(command, db, ".")
	if err != nil {
		err = fmt.Errorf("command '%s' error: %w", command, err)
		return err
	}
	return nil
}
