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

func TestQuery(t *testing.T) {
	// Save and restore original apiCall
	original := apiCall
	defer func() { apiCall = original }()

	// Mock API response
	apiCall = func(url string, result interface{}, _ bool) error {
		res := result.(*Resource)
		res.Results = []Result{
			{Name: "pikachu"},
			{Name: "raichu"},
			{Name: "sandshrew"},
		}
		return nil
	}

	// Now call the query
	res, err := query("pokemon", "chu")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res.Results))
	}
	expected := []string{"pikachu", "raichu"}
	for i, r := range res.Results {
		if r.Name != expected[i] {
			t.Errorf("expected result %s, got %s", expected[i], r.Name)
		}
	}
}
