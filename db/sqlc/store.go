package db

import (
	"context"
	"database/sql"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// create a new Queries struct with the transaction as the db
	q := New(tx)
	err = fn(q)
	if err != nil {
		// if there is an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	ToEntry     Entry    `json:"to_entry"`
	FromEntry   Entry    `json:"from_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, and update the account balance.
// If any of the operations fail, it will rollback the transaction and return an error.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// create a new transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}
		// update the sender's account balance
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: arg.Amount * -1,
		})
		if err != nil {
			return err
		}
		// update the receiver's account balance
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}
