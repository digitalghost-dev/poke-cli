package flags

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

// captureOutput redirects stdout to capture any printed output during a function's execution.
func captureOutput(f func()) string {
	// Save the original stdout
	oldStdout := os.Stdout
	// Create a pipe to capture the output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function, capturing its output
	f()

	// Restore the original stdout and close the writer
	err := w.Close()
	if err != nil {
		return fmt.Sprintf("error closing writer: %v", err)
	}
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestLatestDockerImage(t *testing.T) {
	tests := []struct {
		name        string
		mockRunner  func(name string, args ...string) *exec.Cmd
		expectError bool
		expectText  string
	}{
		{
			name: "success",
			mockRunner: func(name string, args ...string) *exec.Cmd {
				return exec.Command("echo", "v1.0.0\n")
			},
			expectError: false,
			expectText:  "Latest Docker image version: v1.0.0",
		},
		{
			name: "error from command",
			mockRunner: func(name string, args ...string) *exec.Cmd {
				return exec.Command("false") // returns error
			},
			expectError: true,
			expectText:  "Error running command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				latestDockerImage(tt.mockRunner)
			})

			assert.Contains(t, output, tt.expectText)
		})
	}
}

func TestLatestRelease(t *testing.T) {
	githubAPIURL := "https://api.github.com/repos/digitalghost-dev/poke-cli/releases/latest"
	output := captureOutput(func() { latestRelease(githubAPIURL) })

	assert.Contains(t, output, "Latest release tag: v")
}

func TestLatestRelease_Success(t *testing.T) {
	// Create a mock server that simulates a successful response
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, `{"tag_name": "v1.0.0"}`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
			return
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Capture output of the function
	output := captureOutput(func() { latestRelease(server.URL) })
	assert.Contains(t, output, "Latest release tag: v1.0.0")
}

func TestLatestRelease_InvalidJSON(t *testing.T) {
	// Create a mock server that returns invalid JSON
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, `invalid-json`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
			return
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Capture output of the function
	output := captureOutput(func() { latestRelease(server.URL) })
	assert.Contains(t, output, "Error unmarshalling JSON:")
}

func TestLatestFlag(t *testing.T) {
	// Capture the output of the LatestFlag function
	output := captureOutput(LatestFlag)

	// Verify that the output contains expected messages from both latestDockerImage and latestRelease
	assert.Contains(t, output, "Latest Docker image version:", "Expected output to contain 'Latest Docker image version:' but got: %v", output)
	assert.Contains(t, output, "Latest release tag:", "Expected output to contain 'Latest release tag:' but got: %v", output)
}
