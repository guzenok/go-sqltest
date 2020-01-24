package recorder

import (
	"database/sql/driver"
)

type tx struct {
	orig driver.Tx
	rec
}

func (cn *conn) newTx(orig driver.Tx) *tx {
	return &tx{
		orig: orig,
		rec:  cn.rec,
	}
}

// Commit implements driver.Tx.
func (tx *tx) Commit() error {
	after := tx.mock.ExpectCommit()
	tx.write("mock.ExpectCommit()")
	err := tx.orig.Commit()
	if err != nil {
		after.WillReturnError(err)
		tx.write(".WillReturnError(%s)", errToString(err))
	}
	tx.write("\n")
	return err
}

// Rollback implements driver.Tx.
func (tx *tx) Rollback() error {
	after := tx.mock.ExpectRollback()
	tx.write("mock.ExpectRollback()")
	err := tx.orig.Rollback()
	if err != nil {
		after.WillReturnError(err)
		tx.write(".WillReturnError(%s)", errToString(err))
	}
	tx.write("\n")
	return err
}
