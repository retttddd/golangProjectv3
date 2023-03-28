package service

import (
	"awesomeProject3/service/ciphering"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const notSecretVal = "decryptedValue"
const secretVal = "encryptedValue"
const encryptKey = "encryptedKey"
const valFromStorage = "valueFromStorage"
const password = "passwordIs32bytelongandOtherTxt"
const read = "Read"

type mockStorage struct{ mock.Mock }

func newMockStorage() *mockStorage { return &mockStorage{} }

func (m *mockStorage) Read(key string) (string, error) {
	arg := m.Called(key)
	return arg.String(0), arg.Error(1)
}

func (m *mockStorage) Write(key string, value string) error {
	arg := m.Called(key, value)
	return arg.Error(0)
}

type mockEncoder struct{ mock.Mock }

func newMockEncoder() *mockEncoder { return &mockEncoder{} }

func (m *mockEncoder) Encrypt(plaintext string, cipherKey []byte) (string, error) {
	arg := m.Called(plaintext, cipherKey)
	return arg.String(0), arg.Error(1)
}
func (m *mockEncoder) Decrypt(ct string, cipherKey []byte) (string, error) {
	arg := m.Called(ct, cipherKey)
	return arg.String(0), arg.Error(1)
}

func TestSimpleSecretService_WriteSecret(t *testing.T) {
	var testCases = []struct {
		name                 string
		encryptVal           string
		encryptKey           string
		encryptValueErr      error
		encryptKeyErr        error
		value                string
		writeErr             error
		secretKey            string
		key                  string
		numberOfCallsWrite   int
		numberOfCallsEncrypt int
		expectedErr          bool
	}{
		{
			name:                 "Writes data successfully",
			encryptVal:           secretVal,
			secretKey:            password,
			key:                  "key",
			value:                notSecretVal,
			numberOfCallsWrite:   1,
			numberOfCallsEncrypt: 2,
			expectedErr:          false,
			encryptKey:           encryptKey,
		},
		{
			name:                 "Returns error when fails to write data",
			encryptVal:           secretVal,
			secretKey:            password,
			key:                  "key",
			value:                notSecretVal,
			writeErr:             errors.New("write Error"),
			numberOfCallsWrite:   1,
			numberOfCallsEncrypt: 2,
			expectedErr:          true,
			encryptKey:           encryptKey,
		},
		{
			name:                 "Returns error when key is corrupted",
			encryptVal:           secretVal,
			secretKey:            password,
			encryptKeyErr:        errors.New("key Encryption error"),
			key:                  "key",
			value:                notSecretVal,
			numberOfCallsWrite:   0,
			numberOfCallsEncrypt: 2,
			expectedErr:          true,
			encryptKey:           encryptKey,
		},
		{
			name:                 "Returns error when fails to encrypt",
			encryptVal:           secretVal,
			secretKey:            password,
			key:                  "key",
			value:                notSecretVal,
			encryptValueErr:      errors.New("encrypt Error"),
			numberOfCallsWrite:   0,
			numberOfCallsEncrypt: 1,
			expectedErr:          true,
			encryptKey:           encryptKey,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			me := newMockEncoder()
			defer me.AssertNumberOfCalls(t, "Encrypt", tc.numberOfCallsEncrypt)
			me.On("Encrypt", tc.value, ciphering.PassToSecretKey(tc.secretKey)).Return(tc.encryptVal, tc.encryptValueErr)
			me.On("Encrypt", tc.key, ciphering.PassToSecretKey(tc.secretKey)).Return(tc.encryptKey, tc.encryptKeyErr)
			ms := newMockStorage()
			defer ms.AssertNumberOfCalls(t, "Write", tc.numberOfCallsWrite)
			ms.On("Write", tc.encryptKey, tc.encryptVal).Return(tc.writeErr)

			ss := New(ms, me, me)
			err := ss.WriteSecret("key", tc.value, tc.secretKey)
			if tc.expectedErr {
				require.Error(t, err)

			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSimpleSecretService_ReadSecret(t *testing.T) {

	var testCases = []struct {
		name                 string
		readVal              string
		readErr              error
		encryptVal           string
		encryptErr           error
		decryptVal           string
		decryptErr           error
		secretKey            string
		expectedResult       string
		expectedError        bool
		numberOfCallsRead    int
		numberOfCallsEncrypt int
		numberOfCallsDecrypt int
	}{
		{
			name:                 "Reads data successfully",
			readVal:              secretVal,
			encryptVal:           secretVal,
			decryptVal:           notSecretVal,
			secretKey:            password,
			expectedResult:       notSecretVal,
			expectedError:        false,
			numberOfCallsRead:    1,
			numberOfCallsEncrypt: 1,
			numberOfCallsDecrypt: 1,
		},
		{
			name:                 "Returns error when fails to encrypt key",
			readVal:              secretVal,
			encryptVal:           secretVal,
			encryptErr:           errors.New("encryptError"),
			decryptVal:           notSecretVal,
			secretKey:            password,
			expectedResult:       "",
			expectedError:        true,
			numberOfCallsRead:    0,
			numberOfCallsEncrypt: 1,
			numberOfCallsDecrypt: 0,
		},
		{
			name:                 "Returns error when fails to read data",
			readVal:              secretVal,
			readErr:              errors.New("ReadError"),
			encryptVal:           secretVal,
			decryptVal:           notSecretVal,
			secretKey:            password,
			expectedResult:       "",
			expectedError:        true,
			numberOfCallsRead:    1,
			numberOfCallsEncrypt: 1,
			numberOfCallsDecrypt: 0,
		},
		{
			name:                 "Returns error when fails to decrypt read value",
			readVal:              secretVal,
			decryptErr:           errors.New("DecryptError"),
			encryptVal:           secretVal,
			decryptVal:           notSecretVal,
			secretKey:            password,
			expectedResult:       "",
			expectedError:        true,
			numberOfCallsRead:    1,
			numberOfCallsEncrypt: 1,
			numberOfCallsDecrypt: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := newMockStorage()
			defer ms.AssertNumberOfCalls(t, read, tc.numberOfCallsRead)
			ms.On(read, tc.readVal).Return(valFromStorage, tc.readErr)

			me := newMockEncoder()
			defer me.AssertNumberOfCalls(t, "Encrypt", tc.numberOfCallsEncrypt)
			defer me.AssertNumberOfCalls(t, "Decrypt", tc.numberOfCallsDecrypt)
			me.On("Encrypt", tc.decryptVal, ciphering.PassToSecretKey(tc.secretKey)).Return(tc.encryptVal, tc.encryptErr)
			me.On("Decrypt", valFromStorage, ciphering.PassToSecretKey(tc.secretKey)).Return(tc.decryptVal, tc.decryptErr)

			ss := New(ms, me, me)
			result, err := ss.ReadSecret(tc.decryptVal, tc.secretKey)
			if tc.expectedError {
				require.Empty(t, result)
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, result, tc.decryptVal)
			}
		})
	}
}
