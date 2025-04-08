package move

import (
	"bytes"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/styling"
	"log"
	"os"
	"strings"
	"testing"
)

func captureMoveOutput(f func()) string {
	r, w, _ := os.Pipe()
	defer func(r *os.File) {
		err := r.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(r)

	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	f()

	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestMoveCommand(t *testing.T) {
	err := os.Setenv("GO_TESTING", "1")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Unsetenv("GO_TESTING")
		if err != nil {
			fmt.Println(err)
		}
	}()

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "Help flag",
			args:           []string{"move", "--help"},
			expectedOutput: "Get details about a specific move.",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			os.Args = append([]string{"poke-cli"}, tt.args...)

			output := captureMoveOutput(func() {
				defer func() {
					if r := recover(); r != nil {
						if !tt.expectedError {
							t.Fatalf("Unexpected error: %v", r)
						}
					}
				}()
				MoveCommand()
			})

			cleanOutput := styling.StripANSI(output)

			if !strings.Contains(cleanOutput, tt.expectedOutput) {
				t.Errorf("Output mismatch. Expected to contain:\n%s\nGot:\n%s", tt.expectedOutput, output)
			}
		})
	}
}
