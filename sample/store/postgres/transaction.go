package postgres

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type transactionFunc func(tx *sqlx.Tx) error

func (c *connection) ExecuteInTransaction(ctx context.Context, f transactionFunc) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed begin transaction")
	}
	err = f(tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			log.Print(errors.Wrap(errRollback, "failed rollback transaction"))
		}
		return err
	}
	return errors.Wrap(tx.Commit(), "failed commit transaction")
}
