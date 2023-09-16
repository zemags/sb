package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferFx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run a concurrent transfer 100 times
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		// run the transfer in a goroutine
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmptyf(t, result.Transfer.ID, "transfer ID should not be empty")

		// check result
		transfer := result.Transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotEmpty(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.FromAccountID)

		_, err = store.GetEntry(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		toEntry := result.ToEntry
		require.NotEmptyf(t, toEntry, "toEntry should not be empty")
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotEmpty(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmptyf(t, fromEntry, "fromEntry should not be empty")
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotEmpty(t, fromEntry.CreatedAt)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check account balances
		fromAccount := result.FromAccount
		require.NotEmptyf(t, fromAccount, "fromAccount should not be empty")
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmptyf(t, toAccount, "toAccount should not be empty")
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

	}
}
