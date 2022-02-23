package identity

import "github.com/auth0/go-auth0/management"

type UserManager interface {
	Create(u *management.User, opts ...management.RequestOption) error
}
