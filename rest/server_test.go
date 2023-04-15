package rest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
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

func CheckPort() (string, error) {
	for i := 0; i < 20; i++ {

		port := rand.Intn(20000-10000) + 10000
		address := "localhost:" + strconv.Itoa(port)
		log.Println("start port check")
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if conn != nil {
			conn.Close()
		}
		if err != nil {
			log.Println("connection failed port is ok")
			return strconv.Itoa(port), nil

		}
	}
	return "", errors.New("cant find emptey port")
}

func TestSecretRestAPI_Get(t *testing.T) {
	type mockInput struct {
		key      string
		password string
	}
	type mockOutput struct {
		result string
		err    error
	}
	testCases := []struct {
		name          string
		mockInput     mockInput
		mockOutput    mockOutput
		expectedCode  int
		expectedBody  string
		expectedError error
	}{
		{
			name: "successfully runs server",
			mockInput: mockInput{
				key:      "key",
				password: "password",
			},
			mockOutput: mockOutput{
				result: "value",
				err:    nil,
			},
			expectedCode:  http.StatusOK,
			expectedBody:  "value",
			expectedError: nil,
		},
		{
			name: "Returns error when fails to process mocked method",
			mockInput: mockInput{
				key:      "key",
				password: "password",
			},
			mockOutput: mockOutput{
				result: "",
				err:    errors.New("error"),
			},
			expectedCode:  http.StatusNotFound,
			expectedBody:  "error\n",
			expectedError: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.name)
			serverPort, err := CheckPort()
			if err != nil {
				panic("sdfsdf")
			}
			mockSecretService := newMockSecretService()
			mockSecretService.On("ReadSecret", tt.mockInput.key, tt.mockInput.password).Return(tt.mockOutput.result, tt.mockOutput.err)
			srv := NewSecretRestAPI(mockSecretService, serverPort)
			serverCtx, serverCancel := context.WithCancel(context.Background())
			defer serverCancel()
			go func() {
				srv.Start(serverCtx)
			}()
			time.Sleep(1 * time.Second)
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s?getter=key", serverPort), nil)
			req.Header.Set("X-Cipher", tt.mockInput.password)
			client := http.Client{
				Timeout: 5 * time.Second,
			}
			res, err := client.Do(req)
			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			responseData, _ := ioutil.ReadAll(res.Body)
			require.Equal(t, tt.expectedBody, string(responseData))
			require.Equal(t, tt.expectedCode, res.StatusCode)
			log.Println("test case finalized" + tt.name)
		})
	}
}

func TestSecretRestAPI_Post(t *testing.T) {

	var testCases = []struct {
		name          string
		value         string
		key           string
		password      string
		jBody         string
		expectedCode  int
		expectedBody  string
		expectedError error
	}{

		{
			name:          "successfully operates on sent data",
			value:         "value",
			key:           "key",
			password:      "password",
			jBody:         `{"getter" : "key", "value" : "value"}`,
			expectedCode:  http.StatusCreated,
			expectedError: nil,
			expectedBody:  "",
		},
		{
			name:          "Returns error when fails to process corrupted json",
			value:         "value",
			key:           "key",
			password:      "password",
			jBody:         `{sdasdadasdad""||";;;"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: nil,
			expectedBody:  "invalid character",
		},
		{
			name:          "Returns error when fails to process mocked method",
			value:         "value",
			key:           "key",
			password:      "password",
			jBody:         `{"getter" : "key", "value" : "value"}`,
			expectedCode:  http.StatusInternalServerError,
			expectedError: errors.New("Error while performing write method"),
			expectedBody:  "Error while performing write method\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.name)
			serverPort, err := CheckPort()
			if err != nil {
				panic("sdfsdf")
			}
			mockSecretService := newMockSecretService()
			mockSecretService.On("WriteSecret", tt.key, tt.value, tt.password).Return(tt.expectedError)
			srv := NewSecretRestAPI(mockSecretService, serverPort)
			serverCtx, serverCancel := context.WithCancel(context.Background())
			defer serverCancel()
			go func() {
				srv.Start(serverCtx)
			}()
			time.Sleep(1 * time.Second)
			jsonBody := []byte(tt.jBody)
			bodyReader := bytes.NewReader(jsonBody)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s?getter=key", serverPort), bodyReader)
			req.Header.Set("X-Cipher", "password")

			client := http.Client{
				Timeout: 5 * time.Second,
			}
			res, err := client.Do(req)
			responseData, _ := ioutil.ReadAll(res.Body)
			require.Contains(t, string(responseData), tt.expectedBody)
			require.Equal(t, tt.expectedCode, res.StatusCode)
			require.NoError(t, err)

		})
	}
}
