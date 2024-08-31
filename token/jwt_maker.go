package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	secretKey string
}

var ErrInvalidToken = errors.New("token is invalid")
var ErrTokenExpired = errors.New("token has expired")

const minSecretKeyLen = 32

func NewJwtMaker(secretKey string) (Maker, error) {

	if len(secretKey) < 32 {
		return nil, fmt.Errorf("secretKey size should be minimum %d characters", minSecretKeyLen)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}

	payloadBytes, _ := json.Marshal(payload)
	json.Unmarshal(payloadBytes, &claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// fmt.Println(claims)
	// fmt.Println(maker.secretKey)
	// fmt.Println(token)

	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		fmt.Println(err)
		return "", errors.New("token signing failed")
	}

	return tokenStr, nil
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*Payload, error) {

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, keyfunc)
	if err != nil {
		//TODO: handle errs
		return nil, err
	}

	layout := "2006-01-02T15:04:05.9999999-07:00"
	expiredAtTime, err := time.Parse(layout, claims["expired_at"].(string))
	if err != nil {
		return nil, err
	}

	if expiredAtTime.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	fmt.Printf("token.Claims: %v\n", token.Claims)

	payload := &Payload{}

	for k, v := range claims {
		switch k {
		case "id":
			payload.ID, err = uuid.Parse(v.(string))
			if err != nil {
				payload.ID = uuid.Nil
			}
		case "username":
			payload.Username = v.(string)
		case "created_at":
			layout := "2006-01-02T15:04:05.9999999-07:00"
			payload.CreatedAt, err = time.Parse(layout, v.(string))
			if err != nil {
				payload.CreatedAt = time.Time{}
			}
		case "expired_at":
			layout := "2006-01-02T15:04:05.9999999-07:00"
			payload.ExpriredAt, err = time.Parse(layout, v.(string))
			if err != nil {
				payload.ExpriredAt = time.Time{}
			}
		}
	}

	return payload, nil
}
