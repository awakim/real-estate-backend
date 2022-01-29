package token

import (
	"time"

	"github.com/google/uuid"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific userID and duration
	CreateToken(userID uuid.UUID, duration time.Duration) (string, error)
	// CreateRefreshToken creates a new token and returns its tokenID for a given specific userID and duration
	CreateRefreshToken(userID uuid.UUID, duration time.Duration) (string, string, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
	// CreateTokenPair creates fresh access and refresh tokens for the current user
	// If a previous token is included, the previous token is removed from
	// // the redis white list repository
	// CreateTokenPair(userID uuid.UUID, prevTokenID uuid.UUID)
}
