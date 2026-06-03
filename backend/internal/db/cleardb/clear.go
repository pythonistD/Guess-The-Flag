package cleardb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/queries"
)

func ClearDB(db *sqlx.DB) error {
	_, err := db.Exec(queries.CommonQueries.ClearAllTables)
	if err != nil {
		return fmt.Errorf("failed to clear tables: %w", err)
	}
	return nil
}
