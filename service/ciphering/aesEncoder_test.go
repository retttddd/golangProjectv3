package ciphering

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"
const encrpyptedString = "848d56796ac855b0b3f06cb77bb652d2a9a4b84f64026dbfbfe94aaf76cb843608c85a2112ac"
const sourceString = "testString"
const key = "1234567890123456"
const wrongKey = "23"

func TestAesEncoder_Encrypt(t *testing.T) {
	encoder := NewAESEncoder()
	result, err := encoder.Encrypt(sourceString, []byte(key))

	require.Nil(t, err)
	require.Regexp(t, isBase64, result)
}

func TestAesEncoder_Decrypt(t *testing.T) {
	decoder := NewAESEncoder()
	result, err := decoder.Decrypt(encrpyptedString, []byte(key))

	require.Nil(t, err)
	require.Equal(t, sourceString, result)
}

func TestAesEncoder_EncryptWithCorruptedKey(t *testing.T) {
	const WrongKeySizeError = "crypto/aes: invalid key size"

	encoder := NewAESEncoder()
	result, err := encoder.Encrypt(sourceString, []byte(wrongKey))

	require.Empty(t, result)
	require.ErrorContains(t, err, WrongKeySizeError)
}
