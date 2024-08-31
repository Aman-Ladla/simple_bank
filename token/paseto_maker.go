package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symetricKey string) (Maker, error) {

	if len(symetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("length of key should be equals to %d characters", chacha20poly1305.KeySize)
	}

	pasetoMaker := &PasetoMaker{
		paseto:      *paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}
	return pasetoMaker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(tokenStr string) (*Payload, error) {

	payload := &Payload{}

	err := maker.paseto.Decrypt(tokenStr, maker.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if payload.ExpriredAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return payload, nil
}
