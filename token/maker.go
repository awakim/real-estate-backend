package token

import (
	"time"

	"github.com/google/uuid"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific userID and duration
	CreateToken(userID uuid.UUID, duration time.Duration) (*Payload, string, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
	// CreateTokenPair creates fresh access and refresh tokens for the current user
	CreateTokenPair(userID uuid.UUID, accessDuration time.Duration, refreshDuration time.Duration) (Payload, string, Payload, string, error)
}
