package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	pw := RandomString(6)

	hashPw1, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hashPw1)

	err = CheckPassword(pw, hashPw1)
	require.NoError(t, err)

	wrongPw := RandomString(6)
	err = CheckPassword(wrongPw, hashPw1)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashPw2, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hashPw2)
	require.NotEqual(t, hashPw1, hashPw2)
}