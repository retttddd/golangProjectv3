package ciphering

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"
const encryptedString = "848d56796ac855b0b3f06cb77bb652d2a9a4b84f64026dbfbfe94aaf76cb843608c85a2112ac"
const sourceString = "testString"
const key = "1234567890123456"
const wrongKey = "23"

//var encryptionTestCases = []struct {
//	name      string
//	source    string
//	key       string
//	expected  string
//	expectErr bool
//}{
//	{
//		name:      "encryption valid case",
//		source:    sourceString,
//		key:       key,
//		expected:  sourceString,
//		expectErr: false,
//	},
//	{
//		name:      "encryptiom wrong key size",
//		source:    sourceString,
//		key:       wrongKey,
//		expected:  "",
//		expectErr: true,
//	},
//}

func TestAesEncoder_Encrypt_TableDriven(t *testing.T) {
	var encryptionTestCases = []struct {
		name      string
		source    string
		key       string
		expected  string
		expectErr bool
	}{
		{
			name:      "encryption valid case",
			source:    sourceString,
			key:       key,
			expected:  sourceString,
			expectErr: false,
		},
		{
			name:      "encryptiom wrong key size",
			source:    sourceString,
			key:       wrongKey,
			expected:  "",
			expectErr: true,
		},
	}

	for _, tc := range encryptionTestCases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := NewAESEncoder()
			result, err := encoder.Encrypt(tc.source, []byte(tc.key))
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
				require.Regexp(t, isBase64, result)
			}
		})
	}
}

var decryptionTestCases = []struct {
	name      string
	encrypted string
	key       string
	expected  string
	expectErr bool
}{
	{
		name:      "valid case",
		encrypted: encryptedString,
		key:       key,
		expected:  sourceString,
		expectErr: false,
	},
	{
		name:      "wrong key size",
		encrypted: encryptedString,
		key:       "wrongkey",
		expected:  "",
		expectErr: true,
	},
	{
		name:      "corrupted data",
		encrypted: "zzz",
		key:       key,
		expected:  "",
		expectErr: true,
	},
}

func TestAesEncoder_Decrypt_TableDriven(t *testing.T) {
	for _, tc := range decryptionTestCases {
		t.Run(tc.name, func(t *testing.T) {
			decoder := NewAESEncoder()
			result, err := decoder.Decrypt(tc.encrypted, []byte(tc.key))
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}
