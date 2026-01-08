package utils

import (
	"errors"
	"io"
	"os"
	"strings"
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

	output := captureOutput(&os.Stdout, func() {
		HandleCommandOutput(fn)()
	})

	if output != "it worked\n" {
		t.Errorf("expected 'it worked\\n', got %q", output)
	}
}

func TestHandleCommandOutput_Error(t *testing.T) {
	fn := func() (string, error) {
		return "something failed", errors.New("error")
	}

	output := captureOutput(&os.Stderr, func() {
		HandleCommandOutput(fn)()
	})

	if output != "something failed\n" {
		t.Errorf("expected 'something failed\\n', got %q", output)
	}
}

func TestHandleFlagError_WithError(t *testing.T) {
	var b strings.Builder
	msg, err := HandleFlagError(&b, errors.New("bad flag"))

	if got := b.String(); got != "error parsing flags: bad flag\n" {
		t.Fatalf("builder content mismatch: got %q", got)
	}
	if msg != "" {
		t.Fatalf("expected empty message, got %q", msg)
	}
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if err.Error() != "error parsing flags: bad flag" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestHandleFlagError_NilError(t *testing.T) {
	var b strings.Builder
	msg, err := HandleFlagError(&b, nil)

	if got := b.String(); got != "error parsing flags: <nil>\n" {
		t.Fatalf("builder content mismatch for nil error: got %q", got)
	}
	if msg != "" {
		t.Fatalf("expected empty message, got %q", msg)
	}
	if err == nil {
		t.Fatalf("expected non-nil error when wrapping nil")
	}
	// Document current behavior of fmt.Errorf with %w and nil
	if err.Error() != "error parsing flags: %!w(<nil>)" {
		t.Fatalf("unexpected error message for nil wrap: %q", err.Error())
	}
}

func TestWrapText_EmptyString(t *testing.T) {
	if got := WrapText("", 10); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestWrapText_OnlySpaces(t *testing.T) {
	if got := WrapText("   ", 10); got != "   " {
		t.Fatalf("expected to preserve spaces, got %q", got)
	}
}

func TestWrapText_NoWrap(t *testing.T) {
	if got := WrapText("hello world", 20); got != "hello world" {
		t.Fatalf("expected 'hello world', got %q", got)
	}
}

func TestWrapText_CollapseSpaces(t *testing.T) {
	if got := WrapText("hello     world", 20); got != "hello world" {
		t.Fatalf("expected collapsed spaces, got %q", got)
	}
}

func TestWrapText_WrapAtWidth(t *testing.T) {
	if got := WrapText("hello world", 10); got != "hello\nworld" {
		t.Fatalf("expected wrap at width, got %q", got)
	}
}

func TestWrapText_LongWord(t *testing.T) {
	in := "supercalifragilisticexpialidocious"
	if got := WrapText(in, 10); got != in {
		t.Fatalf("expected long word unchanged, got %q", got)
	}
}

func TestWrapText_MultipleLines(t *testing.T) {
	in := "one two three four five"
	expected := "one two\nthree\nfour\nfive"
	if got := WrapText(in, 7); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestCheckHelpFlag_ShortFlag(t *testing.T) {
	// Save and restore original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"poke-cli", "pokemon", "-h"}

	var output strings.Builder
	usageCalled := false
	usageFunc := func() {
		output.WriteString("help message")
		usageCalled = true
	}

	result := CheckHelpFlag(&output, usageFunc)

	if !result {
		t.Error("CheckHelpFlag should return true for -h flag")
	}
	if !usageCalled {
		t.Error("usage function should be called for -h flag")
	}
}

func TestCheckHelpFlag_LongFlag(t *testing.T) {
	// Save and restore original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"poke-cli", "pokemon", "--help"}

	var output strings.Builder
	usageCalled := false
	usageFunc := func() {
		output.WriteString("help message")
		usageCalled = true
	}

	result := CheckHelpFlag(&output, usageFunc)

	if !result {
		t.Error("CheckHelpFlag should return true for --help flag")
	}
	if !usageCalled {
		t.Error("usage function should be called for --help flag")
	}
}

func TestCheckHelpFlag_NoFlag(t *testing.T) {
	// Save and restore original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"poke-cli", "pokemon", "charizard"}

	var output strings.Builder
	usageCalled := false
	usageFunc := func() {
		usageCalled = true
	}

	result := CheckHelpFlag(&output, usageFunc)

	if result {
		t.Error("CheckHelpFlag should return false when no help flag present")
	}
	if usageCalled {
		t.Error("usage function should not be called when no help flag present")
	}
}

func TestCheckHelpFlag_WrongNumberOfArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "too few args",
			args: []string{"poke-cli", "pokemon"},
		},
		{
			name: "too many args",
			args: []string{"poke-cli", "pokemon", "-h", "extra"},
		},
		{
			name: "help flag in wrong position",
			args: []string{"poke-cli", "pokemon", "charizard", "-h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore original os.Args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args

			var output strings.Builder
			usageCalled := false
			usageFunc := func() {
				usageCalled = true
			}

			result := CheckHelpFlag(&output, usageFunc)

			if result {
				t.Errorf("CheckHelpFlag should return false for args: %v", tt.args)
			}
			if usageCalled {
				t.Errorf("usage function should not be called for args: %v", tt.args)
			}
		})
	}
}
