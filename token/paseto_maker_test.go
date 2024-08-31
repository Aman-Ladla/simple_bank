package token

import (
	"testing"
	"time"

	"example.com/simple_bank/db/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {

	jwtMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	createdAt := time.Now()
	expired_at := createdAt.Add(duration)

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, createdAt, payload.CreatedAt, time.Second)
	require.WithinDuration(t, expired_at, payload.ExpriredAt, time.Second)
}

func TestFailPasetoMaker(t *testing.T) {
	jwtMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	_, err = jwtMaker.VerifyToken(token)
	require.Error(t, err, ErrTokenExpired)
}

func TestInvalidTokenWithWrongSymKey(t *testing.T) {

	jwtMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)

	jwtMaker, err = NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	jwtMaker.VerifyToken(token)
	require.Error(t, ErrInvalidToken)
}
