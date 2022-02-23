package identity

import (
	"github.com/auth0/go-auth0/management"
)

func NewUserManager(m *management.Management) *management.UserManager {
	return m.User
}
