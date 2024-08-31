package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	acc1EntArg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    200,
	}

	acc2EntArg := CreateEntryParams{
		AccountID: account2.ID,
		Amount:    -200,
	}

	entry1, err := testQueries.CreateEntry(context.Background(), acc1EntArg)

	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	entry2, err := testQueries.CreateEntry(context.Background(), acc2EntArg)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.Amount, -entry2.Amount)

}
