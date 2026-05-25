package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-bank/internal/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccountID, toAccountID int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount.ID, toAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	transfer := createRandomTransfer(t, fromAccount.ID, toAccount.ID)

	fetchedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedTransfer)

	require.Equal(t, transfer.ID, fetchedTransfer.ID)
	require.Equal(t, transfer.FromAccountID, fetchedTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, fetchedTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, fetchedTransfer.Amount)
	require.WithinDuration(t, transfer.CreatedAt.Time, fetchedTransfer.CreatedAt.Time, time.Second)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for range 10 {
		createRandomTransfer(t, fromAccount.ID, toAccount.ID)
	}

	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}
