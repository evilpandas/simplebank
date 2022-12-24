package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T) Transfer {

	toAccount := createRandomAccount(t)
	fromAccount := createRandomAccount(t)

	args := CreateTransferParams{
		ToAccountID:   toAccount.ID,
		FromAccountID: fromAccount.ID,
		Ammount:       fake.Int64Between(1, 100),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer.ID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.Ammount, transfer.Ammount)
	require.NotEmpty(t, transfer.CreatedAt)

	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.Ammount, transfer2.Ammount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)

}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	args := UpdateTransferParams{
		ID:            transfer1.ID,
		ToAccountID:   transfer1.ToAccountID,
		FromAccountID: transfer1.FromAccountID,
		Ammount:       fake.Int64Between(10, 100),
	}
	transfer2, err := testQueries.UpdateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.NotEqual(t, transfer1.Ammount, transfer2.Ammount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}
	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

}
