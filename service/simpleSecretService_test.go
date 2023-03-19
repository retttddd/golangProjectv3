package service

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const notSecretVal = "decryptedValue"
const secretVal = "encryptedValue"
const secretVal2 = "encryptedValue2"

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

//On("Encrypt", "my-value", mock.AnythingOfType("[]uint8")).Return([]byte("encrypted-value"), nil)
//On("Encrypt", "my-key", mock.AnythingOfType("[]uint8")).Return([]byte("encrypted-key"), nil)
//.On("Write", []byte("encrypted-key"), []byte("encrypted-value")).Return(nil)

func TestSimpleSecretService_WriteSecret(t *testing.T) {
	me := newMockEncoder()
	me.On("Encrypt", notSecretVal, mock.Anything).Return(secretVal, nil)
	me.On("Encrypt", notSecretVal, mock.Anything).Return(secretVal2, nil)

	ms := newMockStorage()
	ms.On("Write", secretVal, secretVal2).Return(nil)

	ss := New(ms, me, me)
	err := ss.WriteSecret("key", notSecretVal, "qweqw")
	require.Nil(t, err)

	//m.AssertNumberOfCalls(t, "Encrypt", 2)

}

func TestSimpleSecretService_ReadSecret(t *testing.T) {
	const read = "Read"
	const valFromStorage = "valueFromStorage"
	ms := newMockStorage()
	ms.On(read, secretVal).Return(valFromStorage, nil)

	me := newMockEncoder()
	me.On("Encrypt", notSecretVal, mock.Anything).Return(secretVal, nil)
	me.On("Decrypt", valFromStorage, mock.Anything).Return(notSecretVal, nil)

	ss := New(ms, me, me)
	result, err := ss.ReadSecret(notSecretVal, "qweqw")
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, result, notSecretVal)
}

var m mockStorage
