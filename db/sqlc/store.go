package db

import (
	"context"
	"database/sql"
	"fmt"
)

// to provide all functions of the DB transactions
type Store struct {
	// composition to embed the queries class in struct
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db), //to create a new queries object
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// now we have queries that run the transaction
	q := New(tx)
	// to get errors of the query
	err = fn(q)

	if err != nil {
		// roleback error
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("Transaction error: %v, Roleback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}

//input parameters
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:amount`
}

//output parameters
type TransferTxResult struct {
	Transfer    Transfer `json:transfer`
	FromAccount Account  `json:from_account`
	ToAccount   Account  `json:to_account`
	FromEntry   Entry    `json:from_entry`
	ToEntry     Entry    `json:to_entry`
}

var txKey = struct{}{} // second empty bracket shows we're creating another empty object of the type

// creates a new transfer record, add account entries, update account balance within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // minus because money is moving out of this account
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // minus because money is moving out of this account
		})

		if err != nil {
			return err
		}

		//get updated accounts' balance

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
