package postgres

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func Connect() (*sql.DB, error) {
	var err error
	dbOnce.Do(func() {
		for i := 0; i < 30; i++ {
			db, err = sql.Open("postgres", "postgres://postgres:postgres@satelite-postgres:5432/satelite?sslmode=disable")
			if err != nil {
				fmt.Printf("Erro ao conectar ao banco de dados: %v\n", err)
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	})
	return db, err
}
