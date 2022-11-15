package terragen_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yardbirdsax/terragen"
	"github.com/yardbirdsax/terragen/internal/terragrunt"
)

func TestDecodeFile(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		expectedOutput terragen.ConfigurationsFile
	}{
		{
			name:     "simple_terragrunt",
			fileName: "simple_terragrunt.hcl",
			expectedOutput: terragen.ConfigurationsFile{
				TerragruntConfigurations: []terragen.TerragruntConfiguration{
					{
						Name:            "test",
						Source:          "mymodule",
						DestinationPath: "path/to/test/terragrunt.hcl",
						Includes: []terragrunt.IncludeConfig{
							{
								Name: "something",
								Path: "hello",
							},
						},
					},
				},
				TerragruntIncludeAlls: []terragrunt.IncludeConfig{
					{
						Name: "all",
						Path: "world",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualOutput := terragen.ConfigurationsFile{}
			filePath := fmt.Sprintf("testdata/%s", tc.fileName)
			err := terragen.DecodeFromFile(filePath, &actualOutput)

			require.NoError(t, err)
			assert.EqualValues(t, tc.expectedOutput, actualOutput)
		})
	}
}

func TestGenerateFromFile(t *testing.T) {
	tests := []struct {
		name                   string
		fileName               string
		expectedOutputFilePath string
		expectedOutput         string
	}{
		{
			name:                   "simple_terragrunt",
			fileName:               "simple_terragrunt.hcl",
			expectedOutputFilePath: "path/to/test/terragrunt.hcl",
			expectedOutput: `
terraform {
  source = "mymodule"
}

include "all" {
  path = "world"
}
include "something" {
  path = "hello"
}
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockFs := afero.NewMemMapFs()
			sourceFilePath := fmt.Sprintf("testdata/%s", tc.fileName)
			err := copyFileToMockFs(sourceFilePath, sourceFilePath, mockFs)
			require.NoError(t, err)
			generator, err := terragen.NewGenerator(terragen.WithFs(mockFs))
			require.NoError(t, err)

			err = generator.GenerateFromFile(sourceFilePath)
			require.NoError(t, err)

			assertFileMatches(t, tc.expectedOutputFilePath, tc.expectedOutput, mockFs)
		})
	}
}

func TestGenerateFromConfig(t *testing.T) {
	tests := []struct {
		name                   string
		config                 *terragen.ConfigurationsFile
		expectedOutputFilePath string
		expectedOutput         string
	}{
		{
			name:                   "simple_terragrunt",
			config: &terragen.ConfigurationsFile{
				TerragruntConfigurations: []terragen.TerragruntConfiguration{
					{
						Name: "test",
						Source: "mymodule",
						DestinationPath: "path/to/test/terragrunt.hcl",
						Includes: []terragrunt.IncludeConfig{
							{
								Name: "something",
								Path: "hello",
							},
						},
					},
				},
				TerragruntIncludeAlls: []terragrunt.IncludeConfig{
					{
						Name: "all",
						Path: "world",
					},
				},
			},
			expectedOutputFilePath: "path/to/test/terragrunt.hcl",
			expectedOutput: `
terraform {
  source = "mymodule"
}

include "all" {
  path = "world"
}
include "something" {
  path = "hello"
}
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockFs := afero.NewMemMapFs()
			generator, err := terragen.NewGenerator(terragen.WithFs(mockFs))
			require.NoError(t, err)

			err = generator.GenerateFromConfig(tc.config)
			require.NoError(t, err)

			assertFileMatches(t, tc.expectedOutputFilePath, tc.expectedOutput, mockFs)
		})
	}
}

func copyFileToMockFs(source string, dest string, mockFs afero.Fs) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := mockFs.Create(dest)
	if err != nil {
		return fmt.Errorf("error opening destination file: %w", err)
	}
	defer destFile.Close()

	var buffer []byte = make([]byte, 1024)
outer:
	for {
		readBytes, err := sourceFile.Read(buffer)
		switch {
		case readBytes == 0:
			break outer
		case err == io.EOF:
			break outer
		case err != nil:
			return fmt.Errorf("error while reading source file: %w", err)
		default:
			_, err = destFile.Write(buffer)
			if err != nil {
				return fmt.Errorf("error while writing to destination file: %w", err)
			}
		}
	}
	return nil
}

func assertFileMatches(t *testing.T, expectedFilePath string, expectedFileContent string, fs afero.Fs) {
	file, err := fs.Open(expectedFilePath)
	require.NoError(t, err)
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.Equal(t, expectedFileContent, string(fileContent))
}
