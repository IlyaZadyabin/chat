package transaction

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"common/database/client"
	"common/database/pg"
)

type manager struct {
	db client.Transactor
}

// NewTransactionManager creates a new transaction manager that satisfies the client.TxManager interface
func NewTransactionManager(db client.Transactor) client.TxManager {
	return &manager{
		db: db,
	}
}

// transaction main function that executes user-specified handler within a transaction
func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn client.Handler) (err error) {
	// If this is a nested transaction, skip starting a new transaction and execute the handler.
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	// Start a new transaction.
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	// Put the transaction in the context.
	ctx = pg.MakeContextTx(ctx, tx)

	// Set up a defer function for transaction rollback or commit.
	defer func() {
		// recover from panic
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		// rollback transaction if an error occurred
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		// if there were no errors, commit the transaction
		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	// Execute code inside the transaction.
	// If the function fails, return an error, and the defer function will perform rollback
	// otherwise the transaction will be committed.
	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f client.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}
