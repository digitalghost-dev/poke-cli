package tcg

import "testing"

func TestFormatInt(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1,000"},
		{4010, "4,010"},
		{10000, "10,000"},
		{1000000, "1,000,000"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatInt(tt.n)
			if got != tt.want {
				t.Errorf("formatInt(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}
