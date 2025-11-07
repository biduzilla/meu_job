package utils

import (
	"database/sql"
	"math/rand"
)

func GenerateRandomCode() int {
	return rand.Intn(900000) + 100000
}

func RunInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
