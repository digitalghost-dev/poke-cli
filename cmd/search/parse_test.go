package search

import "testing"

func TestParseSearch(t *testing.T) {
	mockResults := []Result{
		{Name: "pikachu"},
		{Name: "raichu"},
		{Name: "pidgey"},
		{Name: "sandshrew"},
	}

	tests := []struct {
		name     string
		search   string
		expected []string
	}{
		{
			name:     "Contains match",
			search:   "chu",
			expected: []string{"pikachu", "raichu"},
		},
		{
			name:     "Prefix match",
			search:   "^pi",
			expected: []string{"pikachu", "pidgey"},
		},
		{
			name:     "No match",
			search:   "^z",
			expected: []string{},
		},
		{
			name:     "Contains s",
			search:   "s",
			expected: []string{"sandshrew"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filtered := parseSearch(mockResults, tc.search)
			if len(filtered) != len(tc.expected) {
				t.Errorf("expected %d results, got %d", len(tc.expected), len(filtered))
			}
			for i, result := range filtered {
				if result.Name != tc.expected[i] {
					t.Errorf("expected %s, got %s", tc.expected[i], result.Name)
				}
			}
		})
	}
}
