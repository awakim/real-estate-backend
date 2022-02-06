package db

import (
	"context"
	"testing"

	"github.com/awakim/immoblock-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomUserInfo(t *testing.T) UserInformation {
	user := createRandomUser(t)

	arg := CreateUserInfoParams{
		UserID:      user.ID,
		Firstname:   util.RandomString(6),
		Lastname:    util.RandomString(6),
		PhoneNumber: util.RandomPhoneNumber(),
		Nationality: util.RandomString(6),
		Gender:      "F",
		Address:     util.RandomString(16),
		PostalCode:  util.RandomString(6),
		City:        util.RandomString(6),
		Country:     util.RandomString(6),
	}

	userInfo, err := testQueries.CreateUserInfo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userInfo)

	require.Equal(t, arg.UserID, userInfo.UserID)
	require.Equal(t, arg.Firstname, userInfo.Firstname)
	require.Equal(t, arg.Lastname, userInfo.Lastname)
	require.Equal(t, arg.PhoneNumber, userInfo.PhoneNumber)
	require.Equal(t, arg.Nationality, userInfo.Nationality)
	require.Equal(t, arg.Gender, userInfo.Gender)
	require.Equal(t, arg.Address, userInfo.Address)
	require.Equal(t, arg.PostalCode, userInfo.PostalCode)
	require.Equal(t, arg.City, userInfo.City)
	require.Equal(t, arg.Country, userInfo.Country)

	return userInfo
}

func TestCreateUserInfo(t *testing.T) {
	createRandomUserInfo(t)
}

func TestGetUserInfo(t *testing.T) {
	userInfo := createRandomUserInfo(t)
	userInfo2, err := testQueries.GetUserInfo(context.Background(), userInfo.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, userInfo2)

	require.Equal(t, userInfo2.UserID, userInfo.UserID)
	require.Equal(t, userInfo2.Firstname, userInfo.Firstname)
	require.Equal(t, userInfo2.Lastname, userInfo.Lastname)
	require.Equal(t, userInfo2.PhoneNumber, userInfo.PhoneNumber)
	require.Equal(t, userInfo2.Nationality, userInfo.Nationality)
	require.Equal(t, userInfo2.Gender, userInfo.Gender)
	require.Equal(t, userInfo2.Address, userInfo.Address)
	require.Equal(t, userInfo2.PostalCode, userInfo.PostalCode)
	require.Equal(t, userInfo2.City, userInfo.City)
	require.Equal(t, userInfo2.Country, userInfo.Country)
}
