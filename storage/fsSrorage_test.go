package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		key         string
		fileContent string
		expected    string
		expectedErr error
	}{
		{
			name:        "Key found in the file",
			key:         "foo",
			fileContent: `{"foo": {"value": "bar"}}`,
			expected:    "bar",
			expectedErr: nil,
		},
		{
			name:        "Key not found in the file",
			path:        "",
			key:         "non-existing-key",
			fileContent: `{"foo": {"value": "bar"}}`,
			expected:    "",
			expectedErr: errors.New("item was not found"),
		},
		{
			name:        "Invalid JSON format",
			path:        "",
			key:         "foo",
			fileContent: `invalid-json-format`,
			expected:    "",
			expectedErr: errors.New("invalid character "),
		},
		{
			name:        "File not found",
			path:        "/tmp/non_existing_path/ijrgopiasrais/pSDKF",
			key:         "foo",
			fileContent: "",
			expected:    "",
			expectedErr: errors.New("no such file or directory"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			//using tc.path != "" for test cases which do not need created file
			var targetFilePath string
			if tc.path == "" {

				file, err := ioutil.TempFile("", "prefix")
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(file.Name(), []byte(tc.fileContent), 0644)
				if err != nil {
					log.Fatal(err)
				}
				targetFilePath = file.Name()
			} else {
				targetFilePath = tc.path
			}
			storage := NewFsStorage(targetFilePath)
			result, err := storage.Read(tc.key)
			if tc.expectedErr == nil {
				require.NotNil(t, result)
				require.Equal(t, tc.expected, result)
			} else {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			}

			os.Remove(targetFilePath)
		})

	}

}

func TestWrite(t *testing.T) {

	tests := []struct {
		name        string
		path        string
		key         string
		Data        string
		Value       string
		expected    string
		expectedErr error
	}{
		{
			name:        "valid case",
			path:        "",
			key:         "validKeyNew",
			Data:        `{"validKey": {"value": "validData"}}`,
			Value:       "Value",
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "file was not foundvalid case",
			path:        "/tmp/data.json",
			key:         "validKeyNew",
			Data:        `{"validKey": {"value": "validData"}}`,
			Value:       "Value",
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "file was not found",
			path:        "",
			key:         "validKeyNew",
			Data:        `invalid-json-format`,
			Value:       "Value",
			expected:    "",
			expectedErr: errors.New("invalid character 'i' looking for beginning of value"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var targetFilePath string
			if tc.path == "" {
				file, err := ioutil.TempFile("", "prefix")
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(file.Name(), []byte(tc.Data), 0644)
				if err != nil {
					log.Fatal(err)
				}
				targetFilePath = file.Name()
			} else {
				targetFilePath = tc.path
			}

			storage := NewFsStorage(targetFilePath)
			err := storage.Write(tc.key, tc.Value)

			if tc.expectedErr == nil && err == nil {
				result, errRead := storage.Read(tc.key)
				if errRead != nil {
					require.NoError(t, errRead)
				}
				require.Nil(t, err)
				require.Equal(t, tc.Value, result)
			} else {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			}

			os.Remove(targetFilePath)
		})

	}
}
