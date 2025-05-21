package main

import (
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("valid copy", func(t *testing.T) {
		err := Copy("testdata/input.txt", "test_copy.txt", 0, 0)
		defer os.Remove("test_copy.txt")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("valid copy with limit out of range", func(t *testing.T) {
		err := Copy("testdata/input.txt", "test_copy.txt", 0, 10000)
		defer os.Remove("test_copy.txt")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("invalid offset out of range", func(t *testing.T) {
		err := Copy("testdata/input.txt", "test_copy.txt", 10000, 0)
		defer os.Remove("test_copy.txt")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
