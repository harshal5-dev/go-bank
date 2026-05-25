package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-bank/internal/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)

	fetchedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedEntry)

	require.Equal(t, entry.ID, fetchedEntry.ID)
	require.Equal(t, entry.AccountID, fetchedEntry.AccountID)
	require.Equal(t, entry.Amount, fetchedEntry.Amount)
	require.WithinDuration(t, entry.CreatedAt.Time, fetchedEntry.CreatedAt.Time, time.Second)

}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for range 10 {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

}
