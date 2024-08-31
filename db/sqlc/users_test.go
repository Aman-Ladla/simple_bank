package sqlc

import (
	"context"
	"testing"
	"time"

	"example.com/simple_bank/db/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	userParam := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), userParam)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, userParam.Username, user.Username)
	require.Equal(t, userParam.Email, user.Email)
	require.Equal(t, userParam.FullName, user.FullName)
	require.Equal(t, userParam.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {

	user1 := CreateRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}
