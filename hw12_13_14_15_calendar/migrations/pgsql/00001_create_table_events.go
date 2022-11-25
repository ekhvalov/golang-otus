package migrationspgsql

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up00001, Down00001)
}

func Up00001(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE events (
	id          BIGSERIAL               NOT NULL PRIMARY KEY,
    user_id     BIGINT                  NOT NULL,
    title       CHARACTER VARYING		NOT NULL,
    description text                    NOT NULL,
    start_time  timestamp 				NOT NULL,
    end_time    timestamp   			NOT NULL,
    notify_time timestamp				DEFAULT	'epoch'
);`)
	if err != nil {
		return fmt.Errorf("create table error: %w", err)
	}
	fields := []string{"user_id", "start_time", "end_time", "notify_time"}
	for _, field := range fields {
		_, err = tx.Exec(fmt.Sprintf("CREATE INDEX events_%s_index ON events (%[1]s)", field))
		if err != nil {
			return fmt.Errorf("create index '%s' error: %w", field, err)
		}
	}
	return nil
}

func Down00001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE events;")
	return err
}
