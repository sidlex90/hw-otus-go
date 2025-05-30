package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"VAR1": "value1",
		"VAR2": "value with spaces",
		"VAR3": "",
		"VAR4": "value\x00with\x00nulls",
	}

	for name, content := range testFiles {
		err := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0644)
		require.NoError(t, err)
	}

	// Call ReadDir
	env, err := ReadDir(tempDir)
	require.NoError(t, err)

	// Validate results
	require.Equal(t, "value1", env["VAR1"].Value)
	require.False(t, env["VAR1"].NeedRemove)

	require.Equal(t, "value with spaces", env["VAR2"].Value)
	require.False(t, env["VAR2"].NeedRemove)

	require.Equal(t, "", env["VAR3"].Value)
	require.True(t, env["VAR3"].NeedRemove)

	require.Equal(t, "value\nwith\nnulls", env["VAR4"].Value)
	require.False(t, env["VAR4"].NeedRemove)
}
