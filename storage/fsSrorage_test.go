package storage

import (
	"errors"
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
			name:        "Reads data successfully",
			key:         "foo",
			fileContent: `{"foo": {"value": "bar"}}`,
			expected:    "bar",
			expectedErr: nil,
		},
		{
			name:        "Reads data successfully when key is empty",
			key:         "",
			fileContent: `{"": {"value": "bar"}}`,
			expected:    "bar",
			expectedErr: nil,
		},
		{
			name:        "Reads data successfully when value is empty",
			key:         "key",
			fileContent: `{"key": {"value": ""}}`,
			expected:    "",
			expectedErr: nil,
		},

		{
			name:        "Returns error when fails to find data attached to key",
			path:        "",
			key:         "non-existing-key",
			fileContent: `{"foo": {"value": "bar"}}`,
			expected:    "",
			expectedErr: errors.New("item was not found"),
		},
		{
			name:        "Returns error when fails due to the wrong json format",
			path:        "",
			key:         "foo",
			fileContent: `invalid-json-format`,
			expected:    "",
			expectedErr: errors.New("invalid character "),
		},
		{
			name:        "Returns error when fails to find file",
			path:        "/tmp/non_existing_path/ijrgopiasrais/pSDKF",
			key:         "foo",
			fileContent: "",
			expected:    "",
			expectedErr: errors.New("no such file or directory"),
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
				err = ioutil.WriteFile(file.Name(), []byte(tc.fileContent), 0644)
				if err != nil {
					log.Fatal(err)
				}
				targetFilePath = file.Name()
			} else {
				targetFilePath = tc.path
			}
			defer os.Remove(targetFilePath)
			storage := NewFsStorage(targetFilePath)
			result, err := storage.Read(tc.key)
			if tc.expectedErr == nil {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tc.expected, result)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

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
			name:        "Writes data successfully",
			path:        "",
			key:         "validKeyNew",
			Data:        `{"validKey": {"value": "validData"}}`,
			Value:       "Value",
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "writes data successfully when key is empty",
			path:        "",
			key:         "",
			Data:        `{"": {"value": "validData"}}`,
			Value:       "Value",
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "writes data successfully when value is empty",
			path:        "",
			key:         "keynew",
			Data:        `{"key": {"value": ""}}`,
			Value:       "",
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "writes data successfully when file was not found",
			path:        "/tmp/data.json",
			key:         "validKeyNew",
			Data:        `{"validKey": {"value": "validData"}}`,
			Value:       "Value",
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "Returns error when fails to find and create file",
			path:        "/trere/data.json",
			key:         "validKeyNew",
			Data:        `{"validKey": {"value": "validData"}}`,
			Value:       "Value",
			expected:    "",
			expectedErr: errors.New("no such file or directory"),
		},
		{
			name:        "Returns error when fails to unmarshall data",
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
			defer os.Remove(targetFilePath)
			storage := NewFsStorage(targetFilePath)
			err := storage.Write(tc.key, tc.Value)

			if tc.expectedErr == nil && err == nil {
				result, errRead := storage.Read(tc.key)
				if errRead != nil {
					require.NoError(t, errRead)
				}
				require.NoError(t, err)
				require.Equal(t, tc.Value, result)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

		})

	}
}
