package utils

import (
	"os"
	"reflect"
	"testing"
)

func TestExpandPath(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected string
		hasError bool
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to retrieve home directory: %+v", err)
	}

	tests := []testCase{
		{name: "EmptyPath", input: "", expected: "", hasError: false},
		{name: "RelativePath", input: "documents/test", expected: "documents/test", hasError: false},
		{name: "AbsolutePath", input: "/usr/bin/test", expected: "/usr/bin/test", hasError: false},
		{name: "HomePath", input: "~/Documents", expected: homeDir + "/Documents", hasError: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := ExpandPath(test.input)
			if err != nil && !test.hasError {
				t.Fatalf("Unexpected error: %+v", err)
			}

			if err == nil && test.hasError {
				t.Fatalf("Expected error but got nil")
			}

			if !reflect.DeepEqual(out, test.expected) {
				t.Errorf("Output mismatch. Expected %+v but got %+v", test.expected, out)
			}
		})
	}
}
