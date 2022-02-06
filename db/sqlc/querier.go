// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateProperty(ctx context.Context, arg CreatePropertyParams) (Property, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserInfo(ctx context.Context, arg CreateUserInfoParams) (UserInformation, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	GetProperty(ctx context.Context, id int64) (Property, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
	GetUser(ctx context.Context, email string) (User, error)
	GetUserInfo(ctx context.Context, userID uuid.UUID) (UserInformation, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
