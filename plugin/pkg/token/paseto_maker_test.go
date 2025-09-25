package token

import (
	"crypto/ed25519"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-http-server/core/utils"
	"github.com/stretchr/testify/require"
)

func TestCreatePasetoAndVerifyTokenSuccess(t *testing.T) {
	privateKey := paseto.NewV4AsymmetricSecretKey()
	parser := paseto.NewParserWithoutExpiryCheck()
	maker := NewPasetoMaker(privateKey, parser)

	username := utils.RandomString(6)
	roleId := utils.RandomInt(1, 10)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, roleId, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.Equal(t, username, payload.Username)
	require.Equal(t, roleId, payload.RoleId)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, payload.ExpiredAt, 2*duration)
}

func TestExpiredToken(t *testing.T) {
	privateKey := paseto.NewV4AsymmetricSecretKey()
	parser := paseto.NewParserWithoutExpiryCheck()
	maker := NewPasetoMaker(privateKey, parser)

	token, err := maker.CreateToken(utils.RandomString(6), utils.RandomInt(1, 10), -time.Minute)
	require.NoError(t, err)
	require.NotNil(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, utils.TokenExpired)
	require.Nil(t, payload)
	require.Empty(t, payload)
}

func TestKeyPairMalform(t *testing.T) {
	pub, pri, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	require.NotNil(t, pub)
	require.NotNil(t, pri)

	v4AsymmetricSecretKey, err := paseto.NewV4AsymmetricSecretKeyFromEd25519(pri)
	require.NoError(t, err)
	require.NotNil(t, v4AsymmetricSecretKey)

	parser := paseto.NewParserWithoutExpiryCheck()
	testTokenMaker := NewPasetoMaker(v4AsymmetricSecretKey, parser)
	require.NotEmpty(t, testTokenMaker)

	username := utils.RandomString(6)
	roleId := utils.RandomInt(1, 10)
	duration := time.Minute

	token, err := testTokenMaker.CreateToken(username, roleId, duration)
	require.NoError(t, err)
	require.NotNil(t, token)

	pubServer, priServer, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	require.NotNil(t, pubServer)
	require.NotNil(t, priServer)
	require.NotEqual(t, pub, pubServer)
	require.NotEqual(t, pri, priServer)

	v4AsymmetricSecretKeyServer, err := paseto.NewV4AsymmetricSecretKeyFromEd25519(priServer)
	require.NoError(t, err)
	require.NotNil(t, v4AsymmetricSecretKeyServer)
	require.NotEqual(t, v4AsymmetricSecretKey, v4AsymmetricSecretKeyServer)

	serverTokenMaker := NewPasetoMaker(v4AsymmetricSecretKeyServer, parser)
	payload1, errBadSignature1 := serverTokenMaker.VerifyToken(token)
	require.Error(t, errBadSignature1)
	require.Nil(t, payload1)

	tokenServer, err := serverTokenMaker.CreateToken(username, roleId, duration)
	require.NoError(t, err)
	require.NotEqual(t, token, tokenServer)

	payload2, errBadSignature2 := testTokenMaker.VerifyToken(tokenServer)
	require.Error(t, errBadSignature2)
	require.Nil(t, payload2)

	require.ErrorIs(t, errBadSignature1, errBadSignature2)
}
