package db

import (
	"context"
	"testing"
	"time"

	"github.com/awakim/immoblock-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomProperty(t *testing.T) Property {
	initialBlockCount := util.RandomInt(10, 1000)
	arg := CreatePropertyParams{
		Name:                util.RandomString(6),
		Description:         util.RandomString(32),
		InitialBlockCount:   initialBlockCount,
		RemainingBlockCount: initialBlockCount,
	}

	property, err := testQueries.CreateProperty(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, property)

	require.Equal(t, arg.Name, property.Name)
	require.Equal(t, arg.Description, property.Description)
	require.Equal(t, arg.InitialBlockCount, property.InitialBlockCount)
	require.Equal(t, arg.RemainingBlockCount, property.RemainingBlockCount)

	require.NotZero(t, property.CreatedAt)
	require.NotZero(t, property.UpdatedAt)

	return property
}

func TestCreateProperty(t *testing.T) {
	createRandomProperty(t)
}

func TestGetProperty(t *testing.T) {
	property1 := createRandomProperty(t)
	property2, err := testQueries.GetProperty(context.Background(), property1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, property2)

	require.Equal(t, property1.Name, property2.Name)
	require.Equal(t, property1.Description, property2.Description)
	require.Equal(t, property1.InitialBlockCount, property2.InitialBlockCount)
	require.Equal(t, property1.RemainingBlockCount, property2.RemainingBlockCount)
	require.WithinDuration(t, property1.UpdatedAt, property2.UpdatedAt, time.Second)
	require.WithinDuration(t, property1.CreatedAt, property2.CreatedAt, time.Second)
}
