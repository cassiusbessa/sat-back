package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func CreateDb() error {
	err := createUserTable(db)
	if err != nil {
		return err
	}
	return nil
}

func createUserTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			public_key BYTEA NOT NULL
		)
	`)
	if err != nil {
		println(err.Error())
	}
	return err
}
