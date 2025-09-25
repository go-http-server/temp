package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-http-server/core/utils"
	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	RoleId    int       `json:"role_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	// jwt.RegisteredClaims
}

func NewPayload(username string, roleId int, duration time.Duration) (*Payload, error) {
	uuidv4, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("cannot generate uuid: %s", err)
	}

	payload := &Payload{
		ID:        uuidv4,
		Username:  username,
		RoleId:    roleId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New(utils.TokenExpired)
	}
	return nil
}
