package db

import (
	"context"
	"database/sql"
	"github.com/jaswdr/faker"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

var fake = faker.New()

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    fake.Person().FirstName(),
		Balance:  fake.Int64Between(1, 100),
		Currency: fake.Currency().Code(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("could not connect to db", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
