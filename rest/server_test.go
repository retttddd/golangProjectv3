package rest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type mockSecretService struct{ mock.Mock }

func newMockSecretService() *mockSecretService { return &mockSecretService{} }

func (m *mockSecretService) ReadSecret(key string, password string) (string, error) {
	arg := m.Called(key, password)
	return arg.String(0), arg.Error(1)
}
func (m *mockSecretService) WriteSecret(key string, value string, password string) error {
	arg := m.Called(key, value, password)
	return arg.Error(0)
}

func TestSecretRestAPI_Get_Success(t *testing.T) {
	serverPort := strconv.Itoa(rand.Intn(20000-10000) + 10000)
	mockSecretService := newMockSecretService()
	mockSecretService.On("ReadSecret", "key", "password").Return("value", nil)
	srv := NewSecretRestAPI(mockSecretService, serverPort)
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()
	go func() {
		srv.Start(serverCtx)
	}()
	time.Sleep(1 * time.Second)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s?getter=key", serverPort), nil)
	req.Header.Set("X-Cipher", "password")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	responseData, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "value", string(responseData))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, err)

}
func TestSecretRestAPI_Get_Error(t *testing.T) {
	serverPort := strconv.Itoa(rand.Intn(20000-10000) + 10000)
	mockSecretService := newMockSecretService()
	mockSecretService.On("ReadSecret", "key", "password").Return("", errors.New("error"))
	srv := NewSecretRestAPI(mockSecretService, serverPort)
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()
	go func() {
		srv.Start(serverCtx)
	}()
	time.Sleep(1 * time.Second)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s?getter=key", serverPort), nil)
	req.Header.Set("X-Cipher", "password")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	responseData, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "error\n", string(responseData))
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	require.NoError(t, err)

}
func TestSecretRestAPI_Post_Success(t *testing.T) {
	serverPort := strconv.Itoa(rand.Intn(20000-10000) + 10000)
	mockSecretService := newMockSecretService()
	mockSecretService.On("WriteSecret", "key", "value", "password").Return(nil)
	srv := NewSecretRestAPI(mockSecretService, serverPort)
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()
	go func() {
		srv.Start(serverCtx)
	}()
	time.Sleep(1 * time.Second)
	jsonBody := []byte(`{"getter" : "key", "value" : "value"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s?getter=key", serverPort), bodyReader)
	req.Header.Set("X-Cipher", "password")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	responseData, _ := ioutil.ReadAll(res.Body)
	require.Empty(t, string(responseData))
	require.Equal(t, http.StatusCreated, res.StatusCode) // create all the same
	require.NoError(t, err)

}

func TestSecretRestAPI_Post_FailedToDecodeJSON(t *testing.T) {
	serverPort := strconv.Itoa(rand.Intn(20000-10000) + 10000)
	mockSecretService := newMockSecretService()
	mockSecretService.On("WriteSecret", "key", "value", "password").Return(nil)
	srv := NewSecretRestAPI(mockSecretService, serverPort)
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()
	go func() {
		srv.Start(serverCtx)
	}()
	time.Sleep(1 * time.Second)
	jsonBody := []byte(`{sdasdadasdad""||";;;"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s?getter=key", serverPort), bodyReader)
	req.Header.Set("X-Cipher", "password")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	responseData, _ := ioutil.ReadAll(res.Body)
	require.Contains(t, string(responseData), "invalid character")
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.NoError(t, err)

}
func TestSecretRestAPI_Post_FailedToWriteData(t *testing.T) {
	serverPort := strconv.Itoa(rand.Intn(20000-10000) + 10000)
	mockSecretService := newMockSecretService()
	mockSecretService.On("WriteSecret", "key", "value", "password").Return(errors.New("Some Error"))
	srv := NewSecretRestAPI(mockSecretService, serverPort)
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()
	go func() {
		srv.Start(serverCtx)
	}()
	time.Sleep(1 * time.Second)
	jsonBody := []byte(`{"getter" : "key", "value" : "value"}`)
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s?getter=key", serverPort), bodyReader)
	req.Header.Set("X-Cipher", "password")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	responseData, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, "Some Error\n", string(responseData))
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
	require.NoError(t, err)

}
