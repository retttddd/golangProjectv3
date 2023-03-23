package ciphering

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const regularryEncryptedString = "d4912c64252d453c7776ae2261bce5d11f9dbcdf8ca4d4a92ee1"

func TestAesNotRandomEncoder_Encrypt(t *testing.T) {
	var encryptionTestCases = []struct {
		name   string
		source string

		key       string
		expected  string
		expectErr bool
	}{
		{
			name:      "encryption valid case",
			source:    sourceString,
			key:       key,
			expected:  regularryEncryptedString,
			expectErr: false,
		},
		{
			name:      "encryption wrong key size",
			source:    sourceString,
			key:       wrongKey,
			expected:  "",
			expectErr: true,
		},
	}

	for _, tc := range encryptionTestCases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := NewRegularEncoder()
			result, err := encoder.Encrypt(tc.source, []byte(tc.key))
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestAesNotRandomEncoder_Decrypt(t *testing.T) {
	var decryptionTestCases = []struct {
		name      string
		encrypted string
		key       string
		expected  string
		expectErr bool
	}{
		{
			name:      "valid case",
			encrypted: regularryEncryptedString,
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

	for _, tc := range decryptionTestCases {
		t.Run(tc.name, func(t *testing.T) {
			decoder := NewRegularEncoder()
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
