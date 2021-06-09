package migrations

import (
	"database/sql"
)

func MigrateDatabase(db *sql.DB, queries []string) error {

	var index int
	row := db.QueryRow("SELECT index from migrations")
	_ = row.Scan(&index)

	queries = append([]string{`CREATE TABLE IF NOT EXISTS migrations(index INT);`, `INSERT INTO migrations VALUES(1);`}, queries...)

	for i, query := range queries[index:] {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
		_, err = db.Exec("update migrations set index = $1", i+index+1)
		if err != nil {
			return err
		}
	}

	return nil
}
