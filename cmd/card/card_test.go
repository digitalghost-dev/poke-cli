package card

import (
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
			args:     []string{"card", "-h"},
			wantErr:  false,
			contains: "USAGE:",
		},
		{
			name:     "help flag long",
			args:     []string{"card", "--help"},
			wantErr:  false,
			contains: "FLAGS:",
		},
		{
			name:     "invalid args",
			args:     []string{"card", "invalid-arg"},
			wantErr:  true,
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := CardCommand(tt.args)

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
