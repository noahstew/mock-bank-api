package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	pw := RandomString(6)

	hashPw, err := HashPassword(pw)

	require.NoError(t, err)
	require.NotEmpty(t, hashPw)

	err = CheckPassword(pw, hashPw)
	require.NoError(t, err)

	wrongPw := RandomString(6)
	err = CheckPassword(wrongPw, hashPw)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}