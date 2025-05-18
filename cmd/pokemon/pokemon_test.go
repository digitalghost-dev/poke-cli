package cmd

import (
	"bytes"
	"fmt"
	"github.com/digitalghost-dev/poke-cli/styling"
	"os"
	"strings"
	"testing"
)

func capturePokemonOutput(f func()) string {
	// Create a pipe to capture standard output
	r, w, _ := os.Pipe()
	defer func(r *os.File) {
		err := r.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r)

	// Redirect os.Stdout to the write end of the pipe
	oldStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	// Run the function
	f()

	// Close the write end of the pipe
	err := w.Close()
	if err != nil {
		return ""
	}

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestPokemonCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "Valid abilities flags",
			args:           []string{"pokemon", "sandile", "--abilities"},
			expectedOutput: styling.StripANSI("Your selected Pokémon: Sandile\n• National Pokédex #: 551\n• Weight: 15.2kg (33.5 lbs)\n• Height: 2.3m (2′04″)\n─────────\nAbilities\nAbility 1: Intimidate\nAbility 2: Moxie\nHidden Ability: Anger Point"),
			expectedError:  false,
		},
		{
			name:           "Stats flags",
			args:           []string{"pokemon", "palafin-zero", "--stats"},
			expectedOutput: styling.StripANSI("Your selected Pokémon: Palafin Zero\n• National Pokédex #: 964\n• Weight: 60.2kg (132.7 lbs)\n• Height: 4.3m (4′03″)\n──────────\nBase Stats\nHP         ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 100\nAtk        ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 70\nDef        ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 72\nSp. Atk    ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 53\nSp. Def    ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 62\nSpeed      ▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇ 100\nTotal      457"),
			expectedError:  false,
		},
		{
			name:           "Types flags",
			args:           []string{"pokemon", "armarouge", "--types"},
			expectedOutput: styling.StripANSI("Your selected Pokémon: Armarouge\n• National Pokédex #: 936\n• Weight: 85.0kg (187.4 lbs)\n• Height: 4.9m (4′11″)\n──────\nTyping\nType 1: Fire\nType 2: Psychic"),
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set os.Args for the test
			os.Args = append([]string{"poke-cli"}, tt.args...)

			// Capture the output
			output := capturePokemonOutput(func() {
				defer func() {
					// Recover from os.Exit calls
					if r := recover(); r != nil {
						if !tt.expectedError {
							t.Fatalf("Unexpected error: %v", r)
						}
					}
				}()
				PokemonCommand()
			})

			cleanOutput := styling.StripANSI(output)

			// Check output
			if !strings.Contains(cleanOutput, tt.expectedOutput) {
				t.Logf("DEBUG: Full captured output:\n%s", cleanOutput)
				t.Errorf("Output mismatch.\nExpected to contain:\n%s\nGot:\n%s", tt.expectedOutput, cleanOutput)
			}
		})
	}
}
