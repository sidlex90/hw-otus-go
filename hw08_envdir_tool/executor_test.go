package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Test case: Run a valid command with environment variables
	t.Run("Run valid command with environment variables", func(t *testing.T) {
		outputFileName := "/tmp/envOutput-1.txt"

		defer os.Remove(outputFileName)

		env := Environment{
			"TEST_VAR": {Value: "test_value", NeedRemove: false},
		}
		cmd := []string{"./testdata/envToFile.sh", outputFileName}

		returnCode := RunCmd(cmd, env)
		output, err := os.ReadFile(outputFileName)
		require.NoError(t, err)
		require.Equal(t, 0, returnCode)
		require.Contains(t, string(output), "TEST_VAR=test_value")
	})

	// Test case: Remove an environment variable
	t.Run("Remove environment variable", func(t *testing.T) {
		outputFileName := "/tmp/envOutput-2.txt"

		env := Environment{
			"TEST_VAR": {Value: "", NeedRemove: true},
		}
		cmd := []string{"./testdata/envToFile.sh", outputFileName}

		// Set the environment variable before running the command
		os.Setenv("TEST_VAR", "should_be_removed")

		returnCode := RunCmd(cmd, env)
		defer os.Remove(outputFileName)

		// Restore stdout and read the output
		output, err := os.ReadFile(outputFileName)
		require.NoError(t, err)
		require.Equal(t, 0, returnCode)
		require.NotContains(t, string(output), "TEST_VAR=")
	})
}
