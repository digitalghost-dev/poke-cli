package card

import (
	"os"
	"strings"
	"testing"
)

func TestCardCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "help flag short",
			args:     []string{"poke-cli", "card", "-h"},
			wantErr:  false,
			contains: "USAGE:",
		},
		{
			name:     "help flag long",
			args:     []string{"poke-cli", "card", "--help"},
			wantErr:  false,
			contains: "FLAGS:",
		},
		{
			name:     "invalid args",
			args:     []string{"poke-cli", "card", "invalid-arg"},
			wantErr:  true,
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			os.Args = tt.args
			defer func() { os.Args = oldArgs }()

			output, err := CardCommand()

			if (err != nil) != tt.wantErr {
				t.Errorf("CardCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Errorf("CardCommand() output should contain %q, got %q", tt.contains, output)
			}
		})
	}
}
