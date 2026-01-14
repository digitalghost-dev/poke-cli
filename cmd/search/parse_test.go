package search

import "testing"

func TestParseSearch(t *testing.T) {
	mockResults := []Result{
		{Name: "hariyama"},
		{Name: "pikachu"},
		{Name: "raichu"},
		{Name: "pidgey"},
		{Name: "sandshrew"},
		{Name: "bulbasaur"},
		{Name: "charmander"},
		{Name: "charmeleon"},
		{Name: "squirtle"},
		{Name: "musharna"},
		{Name: "caterpie"},
		{Name: "weedle"},
		{Name: "rattata"},
	}

	t.Run("Substring match prefers contains over fuzzy", func(t *testing.T) {
		filtered, err := parseSearch(mockResults, "hari")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(filtered) != 1 {
			t.Fatalf("expected 1 result, got %d", len(filtered))
		}
		if filtered[0].Name != "hariyama" {
			t.Fatalf("expected hariyama, got %s", filtered[0].Name)
		}
	})

	t.Run("Regex prefix match ^", func(t *testing.T) {
		filtered, err := parseSearch(mockResults, "^pi")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"pikachu", "pidgey"}
		if len(filtered) != len(expected) {
			t.Fatalf("expected %d results, got %d", len(expected), len(filtered))
		}
		for i, result := range filtered {
			if result.Name != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], result.Name)
			}
		}
	})

	t.Run("Regex suffix match $", func(t *testing.T) {
		filtered, err := parseSearch(mockResults, "chu$")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"pikachu", "raichu"}
		if len(filtered) != len(expected) {
			t.Fatalf("expected %d results, got %d", len(expected), len(filtered))
		}
		for i, result := range filtered {
			if result.Name != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], result.Name)
			}
		}
	})

	t.Run("Regex no match", func(t *testing.T) {
		filtered, err := parseSearch(mockResults, "^z")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(filtered) != 0 {
			t.Errorf("expected 0 results, got %d", len(filtered))
		}
	})

	t.Run("Regex character class []", func(t *testing.T) {
		filtered, err := parseSearch(mockResults, "^[rs]")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"raichu", "sandshrew", "squirtle", "rattata"}
		if len(filtered) != len(expected) {
			t.Fatalf("expected %d results, got %d", len(expected), len(filtered))
		}
		for i, result := range filtered {
			if result.Name != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], result.Name)
			}
		}
	})

	t.Run("Regex alternation |", func(t *testing.T) {
		filtered, err := parseSearch(mockResults, "^(pikachu|bulbasaur)$")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"pikachu", "bulbasaur"}
		if len(filtered) != len(expected) {
			t.Fatalf("expected %d results, got %d", len(expected), len(filtered))
		}
		for i, result := range filtered {
			if result.Name != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], result.Name)
			}
		}
	})

	t.Run("Regex one or more +", func(t *testing.T) {
		// Matches any name containing one or more 'e'
		filtered, err := parseSearch(mockResults, "e+")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"pidgey", "sandshrew", "charmander", "charmeleon", "squirtle", "caterpie", "weedle"}
		if len(filtered) != len(expected) {
			t.Fatalf("expected %d results, got %d", len(expected), len(filtered))
		}
		for i, result := range filtered {
			if result.Name != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], result.Name)
			}
		}
	})

	t.Run("Regex optional ?", func(t *testing.T) {
		filtered, err := parseSearch(mockResults, "^chu?arm")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"charmander", "charmeleon"}
		if len(filtered) != len(expected) {
			t.Fatalf("expected %d results, got %d", len(expected), len(filtered))
		}
		for i, result := range filtered {
			if result.Name != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], result.Name)
			}
		}
	})

	t.Run("Invalid regex returns error", func(t *testing.T) {
		_, err := parseSearch(mockResults, "[invalid")
		if err == nil {
			t.Error("expected error for invalid regex, got nil")
		}
	})
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

	// Now call the query with regex pattern
	res, err := query("pokemon", ".*chu.*")
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
