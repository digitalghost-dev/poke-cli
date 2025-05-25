package utils

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func captureOutput(target **os.File, fn func()) string {
	// Create pipe
	r, w, _ := os.Pipe()
	orig := *target
	*target = w

	// Run the function
	fn()

	// Restore original and read output
	err := w.Close()
	if err != nil {
		return ""
	}
	*target = orig
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestHandleCommandOutput_Success(t *testing.T) {
	fn := func() (string, error) {
		return "it worked", nil
	}

	output := captureOutput(&os.Stdout, HandleCommandOutput(fn))

	if output != "it worked\n" {
		t.Errorf("expected 'it worked\\n', got %q", output)
	}
}

func TestHandleCommandOutput_Error(t *testing.T) {
	fn := func() (string, error) {
		return "something failed", fmt.Errorf("error")
	}

	output := captureOutput(&os.Stderr, HandleCommandOutput(fn))

	if output != "something failed\n" {
		t.Errorf("expected 'something failed\\n', got %q", output)
	}
}
