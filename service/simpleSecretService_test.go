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

func TestSimpleSecretService_WriteSecret2(t *testing.T) {
	var testCases = []struct {
		name           string
		encryptVal     string
		encryptErr     error
		value          string
		writeErr       error
		secretKey      string
		key            string
		numberOfCallsW int
		numberOfCallsE int
		expectedErr    bool
	}{
		{
			name:           "Successful Write Secret",
			encryptVal:     secretVal,
			secretKey:      password,
			key:            "key",
			value:          notSecretVal,
			numberOfCallsW: 1,
			numberOfCallsE: 2,
			expectedErr:    false,
		},
		{
			name:           "Unsuccessful encrypt Write Secret",
			encryptVal:     secretVal,
			secretKey:      password,
			key:            "key",
			value:          notSecretVal,
			encryptErr:     errors.New("encrypt Error"),
			numberOfCallsW: 0,
			numberOfCallsE: 1,
			expectedErr:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			me := newMockEncoder()
			defer me.AssertNumberOfCalls(t, "Encrypt", tc.numberOfCallsE)
			me.On("Encrypt", tc.value, ciphering.PassToSecretKey(tc.secretKey)).Return(tc.encryptVal, tc.encryptErr)
			me.On("Encrypt", tc.key, ciphering.PassToSecretKey(tc.secretKey)).Return(tc.encryptVal, tc.encryptErr)

			ms := newMockStorage()
			defer ms.AssertNumberOfCalls(t, "Write", tc.numberOfCallsW)
			ms.On("Write", tc.encryptVal, tc.encryptVal).Return(tc.writeErr)

			ss := New(ms, me, me)
			err := ss.WriteSecret("key", tc.value, tc.secretKey)
			if tc.expectedErr {
				require.Error(t, err)

			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestSimpleSecretService_ReadSecret(t *testing.T) {

	var testCases = []struct {
		name           string
		readVal        string
		readErr        error
		encryptVal     string
		encryptErr     error
		decryptVal     string
		decryptErr     error
		secretKey      string
		key            string
		value          string
		expectedResult string
		expectedError  bool
		numberOfCallsR int
		numberOfCallsE int
		numberOfCallsD int
	}{
		{
			name:           "Successful Read Secret",
			readVal:        secretVal,
			encryptVal:     secretVal,
			decryptVal:     notSecretVal,
			secretKey:      password,
			key:            "key",
			value:          notSecretVal,
			expectedResult: notSecretVal,
			expectedError:  false,
			numberOfCallsR: 1,
			numberOfCallsE: 1,
			numberOfCallsD: 1,
		},
		{
			name:           "Read Secret Encrypt Error",
			readVal:        secretVal,
			encryptVal:     secretVal,
			encryptErr:     errors.New("encryptError"),
			decryptVal:     notSecretVal,
			secretKey:      password,
			key:            "key",
			value:          notSecretVal,
			expectedResult: "",
			expectedError:  true,
			numberOfCallsR: 0,
			numberOfCallsE: 1,
			numberOfCallsD: 0,
		},
		{
			name:           "Read Secret read Error",
			readVal:        secretVal,
			readErr:        errors.New("ReadError"),
			encryptVal:     secretVal,
			decryptVal:     notSecretVal,
			secretKey:      password,
			key:            "key",
			value:          notSecretVal,
			expectedResult: "",
			expectedError:  true,
			numberOfCallsR: 1,
			numberOfCallsE: 1,
			numberOfCallsD: 0,
		},
		{
			name:           "Read Secret Decrypt Error",
			readVal:        secretVal,
			decryptErr:     errors.New("DecryptError"),
			encryptVal:     secretVal,
			decryptVal:     notSecretVal,
			secretKey:      password,
			key:            "key",
			value:          notSecretVal,
			expectedResult: "",
			expectedError:  true,
			numberOfCallsR: 1,
			numberOfCallsE: 1,
			numberOfCallsD: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms := newMockStorage()
			defer ms.AssertNumberOfCalls(t, read, tc.numberOfCallsR)
			ms.On(read, tc.readVal).Return(valFromStorage, tc.readErr)

			me := newMockEncoder()
			defer me.AssertNumberOfCalls(t, "Encrypt", tc.numberOfCallsE)
			defer me.AssertNumberOfCalls(t, "Decrypt", tc.numberOfCallsD)
			me.On("Encrypt", notSecretVal, ciphering.PassToSecretKey(password)).Return(secretVal, tc.encryptErr)
			me.On("Decrypt", valFromStorage, ciphering.PassToSecretKey(password)).Return(notSecretVal, tc.decryptErr)

			ss := New(ms, me, me)
			result, err := ss.ReadSecret(notSecretVal, password)
			if tc.expectedError {
				require.Empty(t, result)
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, result, notSecretVal)
			}
		})
	}
}

//func TestSimpleSecretService_WriteSecret(t *testing.T) {
//
//	me := newMockEncoder()
//	defer me.AssertExpectations(t)
//	me.On("Encrypt", notSecretVal, ciphering.PassToSecretKey(password)).Return(secretVal, nil)
//	me.On("Encrypt", "key", ciphering.PassToSecretKey(password)).Return(secretVal, nil)
//
//	ms := newMockStorage()
//	defer ms.AssertExpectations(t)
//	ms.On("Write", secretVal, secretVal).Return(nil)
//
//	ss := New(ms, me, me)
//	err := ss.WriteSecret("key", notSecretVal, password)
//	require.Nil(t, err)
//
//}
//
//func TestSimpleSecretService_WriteSecret_Errors(t *testing.T) {
//	ms := newMockStorage()
//	me := newMockEncoder()
//	ss := New(ms, me, me)
//
//	defer me.AssertExpectations(t)
//
//	me.On("Encrypt", notSecretVal, ciphering.PassToSecretKey(password)).Return("", errors.New("encryption error"))
//	err := ss.WriteSecret("key", notSecretVal, password)
//	assert.Error(t, err)
//
//}

//func TestSimpleSecretService_ReadSecret(t *testing.T) {
//
//	ms := newMockStorage()
//	ms.On(read, secretVal).Return(valFromStorage, nil)
//
//	me := newMockEncoder()
//	me.On("Encrypt", notSecretVal, ciphering.PassToSecretKey(password)).Return(secretVal, nil)
//	me.On("Decrypt", valFromStorage, ciphering.PassToSecretKey(password)).Return(notSecretVal, nil)
//
//	ss := New(ms, me, me)
//	result, err := ss.ReadSecret(notSecretVal, password)
//	require.Nil(t, err)
//	require.NotNil(t, result)
//	require.Equal(t, result, notSecretVal)
//}
//
//func TestSimpleSecretService_ReadSecret_EncryptError(t *testing.T) {
//	ms := newMockStorage()
//	me := newMockEncoder()
//	ss := New(ms, me, me)
//
//	defer me.AssertExpectations(t)
//
//	me.On("Encrypt", "key", ciphering.PassToSecretKey(password)).Return("", errors.New("key encryption error"))
//	result, err := ss.ReadSecret("key", password)
//	require.Empty(t, result)
//	assert.Error(t, err)
//}
//
//func TestSimpleSecretService_ReadSecret_ReadError(t *testing.T) {
//	ms := newMockStorage()
//	me := newMockEncoder()
//	ss := New(ms, me, me)
//
//	defer me.AssertExpectations(t)
//	me.On("Encrypt", "key", ciphering.PassToSecretKey(password)).Return(secretVal2, nil)
//
//	defer ms.AssertExpectations(t)
//	ms.On(read, secretVal2).Return("", errors.New("reading file from storage error"))
//
//	result, err := ss.ReadSecret("key", password)
//	require.Empty(t, result)
//	assert.Error(t, err)
//}
//func TestSimpleSecretService_ReadSecret_DecryptError(t *testing.T) {
//	ms := newMockStorage()
//	me := newMockEncoder()
//	ss := New(ms, me, me)
//
//	defer me.AssertExpectations(t)
//	me.On("Encrypt", "key", ciphering.PassToSecretKey(password)).Return(secretVal, nil)
//	ms.On(read, secretVal).Return(valFromStorage, nil)
//
//	defer ms.AssertExpectations(t)
//	me.On("Decrypt", valFromStorage, ciphering.PassToSecretKey(password)).Return("", errors.New("Encrypt Error"))
//
//	result, err := ss.ReadSecret("key", password)
//	require.Empty(t, result)
//	assert.Error(t, err)
//}
