package token

import (
	"time"

	"github.com/google/uuid"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific userID and duration
	CreateToken(userID uuid.UUID, duration time.Duration) (string, string, error)
	// CreateRefreshToken creates a new token and returns its tokenID for a given specific userID and duration
	CreateRefreshToken(userID uuid.UUID, duration time.Duration) (string, string, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
	// CreateTokenPair creates fresh access and refresh tokens for the current user
	CreateTokenPair(userID uuid.UUID, accessDuration time.Duration, refreshDuration time.Duration) (string, string, string, string, error)
}
