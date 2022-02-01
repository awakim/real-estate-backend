package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// CreateToken creates a new token for a specific userID and duration
func (maker *PasetoMaker) CreateToken(userID uuid.UUID, duration time.Duration) (*Payload, string, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return &Payload{}, "", err
	}

	st, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return payload, st, err
}

// CreateTokenPair creates a new pair of tokens for a specific userID and duration
func (maker *PasetoMaker) CreateTokenPair(userID uuid.UUID, accessDuration time.Duration, refreshDuration time.Duration) (Payload, string, Payload, string, error) {
	at, atST, err := maker.CreateToken(
		userID,
		accessDuration,
	)
	if err != nil {
		return Payload{}, "", Payload{}, "", err
	}

	rt, rtST, err := maker.CreateToken(
		userID,
		refreshDuration,
	)
	if err != nil {
		return Payload{}, "", Payload{}, "", err
	}
	return *at, atST, *rt, rtST, nil
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
