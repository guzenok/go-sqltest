package recorder

import (
	"database/sql/driver"
)

type tx struct {
	orig driver.Tx
	*output
}

func (cn *conn) newTx(orig driver.Tx) *tx {
	return &tx{
		orig:   orig,
		output: cn.output,
	}
}

// Commit implements driver.Tx.
func (tx *tx) Commit() error {
	after := tx.mock.ExpectCommit()
	tx.p("mock.ExpectCommit()")
	err := tx.orig.Commit()
	if err != nil {
		after.WillReturnError(err)
		tx.p(".WillReturnError(%s)", tx.errToString(err))
	}
	tx.p("\n")
	return err
}

// Rollback implements driver.Tx.
func (tx *tx) Rollback() error {
	after := tx.mock.ExpectRollback()
	tx.p("mock.ExpectRollback()")
	err := tx.orig.Rollback()
	if err != nil {
		after.WillReturnError(err)
		tx.p(".WillReturnError(%s)", tx.errToString(err))
	}
	tx.p("\n")
	return err
}
