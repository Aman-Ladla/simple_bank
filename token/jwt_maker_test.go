package token

import (
	"encoding/json"
	"testing"
	"time"

	"example.com/simple_bank/db/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJwtMaker(t *testing.T) {

	jwtMaker, err := NewJwtMaker(util.RandomString(32))
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

func TestFailJwtMaker(t *testing.T) {
	jwtMaker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	_, err = jwtMaker.VerifyToken(token)
	require.Error(t, err, ErrTokenExpired)
}

func TestInvalidToken(t *testing.T) {

	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	claims := jwt.MapClaims{}

	payloadBytes, _ := json.Marshal(payload)
	json.Unmarshal(payloadBytes, &claims)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)

	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(tokenStr)
	require.Error(t, err)
	require.Empty(t, payload)
}
