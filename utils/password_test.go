package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateHashPassword(t *testing.T) {
	rawPassword := RandomString(6)

	hash1, err := HashPassword(rawPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)
	err = ComparePassword(rawPassword, hash1)
	require.NoError(t, err)

	hash2, err := HashPassword(rawPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)
	err = ComparePassword(rawPassword, hash2)
	require.NoError(t, err)

	require.NotEqual(t, hash1, hash2)

	wrongPassword := RandomString(6)
	err = ComparePassword(wrongPassword, hash1)
	require.Error(t, err)
	err = ComparePassword(wrongPassword, hash2)
	require.Error(t, err)
}
