package migrationspgsql

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up00001, Down00001)
}

func Up00001(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE events (
	id          BIGSERIAL                NOT NULL PRIMARY KEY,
    user_id     BIGINT                   NOT NULL,
    title       CHARACTER VARYING        NOT NULL,
    description text                     NOT NULL,
    start_time  timestamp 				 NOT NULL,
    end_time    timestamp                NOT NULL
);`)
	return err
}

func Down00001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE events;")
	return err
}
