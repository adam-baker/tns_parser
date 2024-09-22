package tnsparser

import (
  "os"
  "path/filepath"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestParseValidConfigurations(t *testing.T) {
    testFiles, err := filepath.Glob("testdata/valid/*.tns")
    if err != nil {
        t.Fatalf("Failed to list test files: %v", err)
    }

    for _, testFile := range testFiles {
        testFile := testFile // Capture range variable
        t.Run(filepath.Base(testFile), func(t *testing.T) {
            content, err := os.ReadFile(testFile)
            assert.NoError(t, err, "Failed to read test file: %s", testFile)

            tnsFile, err := ParseTNSString(string(content))
            assert.NoError(t, err, "Parsing valid TNS file should not produce an error: %s", testFile)
            assert.NotNil(t, tnsFile, "Parsed TNSFile should not be nil")
            assert.NotEmpty(t, tnsFile.Entries, "Parsed TNSFile should have entries")

            // You can perform further assertions based on expected data
            // For example, check the service name
            expectedServiceName := filepath.Base(testFile) // or derive from test file name
            entry := tnsFile.Entries[0]
            assert.Equal(t, expectedServiceName, entry.Name, "Service name should match file name")
        })
    }
}

func TestParseInvalidSyntax(t *testing.T) {
    testFiles, err := filepath.Glob("testdata/invalid_syntax/*.tns")
    if err != nil {
        t.Fatalf("Failed to list test files: %v", err)
    }

    for _, testFile := range testFiles {
        testFile := testFile
        t.Run(filepath.Base(testFile), func(t *testing.T) {
            content, err := os.ReadFile(testFile)
            assert.NoError(t, err, "Failed to read test file: %s", testFile)

            _, err = ParseTNSString(string(content))
            assert.Error(t, err, "Expected an error when parsing invalid syntax: %s", testFile)
            assert.Contains(t, err.Error(), "expected '=' after key", "Error message should indicate missing '='")
        })
    }
}
