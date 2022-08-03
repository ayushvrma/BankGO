package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("before balance: ", account1.Balance, account2.Balance)

	// to make sure our transactions work well, we run them with several go routines
	n := 5
	amount := int64(10)

	// two channels to handle go routines, one to handle errors, the other to handle result
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			// to send data back to the channel
			// channel on left when sending value to it
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	// now check for errors
	for i := 0; i < n; i++ {
		// channel is to the right when recieveing value from it
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)

		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// checking account balance
		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//checkacc balance
		fmt.Println("tx balance: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance //money flowing out of account 1
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)

		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //amount, 2*amount, 3*amount, 4*amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n) //n is number of transaction
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAcc1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc1)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc2)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAcc1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAcc2.Balance)

	fmt.Println("after balance: ", updatedAcc1.Balance, updatedAcc2.Balance)

}
