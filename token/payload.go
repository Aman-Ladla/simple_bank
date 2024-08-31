package token

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	ExpriredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {

	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:         uuid,
		Username:   username,
		CreatedAt:  time.Now(),
		ExpriredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// func (payload *Payload)
