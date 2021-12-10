package boilergql

import (
	"context"
	"database/sql"
	"fmt"
)

func RunInTransaction(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return fmt.Errorf("rollbackErr: %s, err: %s", err, rollBackErr)
		}
		return err
	}

	return tx.Commit()
}


