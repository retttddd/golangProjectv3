package ciphering

import (
	"crypto/rand"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"
const encryptedString = "303030303030303030303030db526531bd980cf6eb3feb4573baef490f65a143517ded1e5b1c"
const sourceString = "testString"
const key = "1234567890123456"
const wrongKey = "23"

type mockNonceProducer struct{ mock.Mock }

func newMockNonceProducer() *mockNonceProducer { return &mockNonceProducer{} }

func (m *mockNonceProducer) generate(size int) (string, error) {
	arg := m.Called(size)
	return arg.String(0), arg.Error(1)
}

func TestAesEncoder_Encrypt(t *testing.T) {
	const size = 12
	mockSuccessNonce := newMockNonceProducer()
	mockSuccessNonce.On("generate", size).Return("000000000000", nil)
	mockFailNonce := newMockNonceProducer()
	mockFailNonce.On("generate", size).Return("", errors.New("Generation fail"))

	var encryptionTestCases = []struct {
		name      string
		source    string
		key       string
		expected  string
		expectErr bool
		nonceMock nonceProducer
	}{
		{
			name:      "Encrypts data successfully",
			source:    sourceString,
			key:       key,
			expected:  encryptedString,
			expectErr: false,
			nonceMock: mockSuccessNonce,
		},
		{
			name:      "Returns error when fails due to wrong key size",
			source:    sourceString,
			key:       wrongKey,
			expected:  "",
			expectErr: true,
			nonceMock: mockSuccessNonce,
		},
		{
			name:      "Returns error when	nonce generation fails",
			source:    sourceString,
			key:       key,
			expected:  "",
			expectErr: true,
			nonceMock: mockFailNonce,
		},
	}

	for _, tc := range encryptionTestCases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := NewAESEncoder(tc.nonceMock)
			result, err := encoder.Encrypt(tc.source, []byte(tc.key))
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Regexp(t, isBase64, result)
				require.Equal(t, encryptedString, result)
			}
		})
	}
}

func TestAesEncoder_Decrypt(t *testing.T) {
	var decryptionTestCases = []struct {
		name      string
		encrypted string
		key       string
		expected  string
		expectErr bool
	}{
		{
			name:      "Decrypts data successfully",
			encrypted: encryptedString,
			key:       key,
			expected:  sourceString,
			expectErr: false,
		},
		{
			name:      "Returns error when fails due to wrong key size",
			encrypted: encryptedString,
			key:       "wrongkey",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Returns error when fails due to corrupted data",
			encrypted: "zzz",
			key:       key,
			expected:  "",
			expectErr: true,
		},
	}

	for _, tc := range decryptionTestCases {
		t.Run(tc.name, func(t *testing.T) {
			decoder := NewAESEncoder(NewRandomNonceProducer(rand.Reader))
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
