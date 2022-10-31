package pgsqlstorage

import "fmt"

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (c Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}
