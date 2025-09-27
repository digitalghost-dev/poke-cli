package connections

import (
	"os"
	"strings"
	"testing"
)

func TestQueryBerryData(t *testing.T) {
	t.Run("basic query without parameters", func(t *testing.T) {
		query := "SELECT name FROM berries LIMIT 1"
		results, err := QueryBerryData(query)

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		if len(results) == 0 {
			t.Error("QueryBerryData() should return at least one result")
		}
	})

	t.Run("query with parameters", func(t *testing.T) {
		query := "SELECT name FROM berries WHERE name LIKE ? LIMIT 1"
		results, err := QueryBerryData(query, "%a%")

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		// Should return some results since many berry names contain 'a'
		if len(results) == 0 {
			t.Error("QueryBerryData() should return at least one result for berries containing 'a'")
		}
	})

	t.Run("query with multiple parameters", func(t *testing.T) {
		query := "SELECT name FROM berries WHERE name LIKE ? OR name LIKE ? LIMIT 5"
		results, err := QueryBerryData(query, "%a%", "%e%")

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		if len(results) == 0 {
			t.Error("QueryBerryData() should return results for berries containing 'a' or 'e'")
		}
	})

	t.Run("invalid query", func(t *testing.T) {
		query := "SELECT invalid_column FROM non_existent_table"
		_, err := QueryBerryData(query)

		if err == nil {
			t.Error("QueryBerryData() should return an error for invalid query")
		}
	})

	t.Run("empty query", func(t *testing.T) {
		query := ""
		_, err := QueryBerryData(query)

		t.Logf("Empty query result: error = %v", err)
	})

	t.Run("query with no results", func(t *testing.T) {
		query := "SELECT name FROM berries WHERE name = ?"
		results, err := QueryBerryData(query, "NonExistentBerryName12345")

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		if len(results) != 0 {
			t.Errorf("QueryBerryData() should return empty results for non-existent berry, got %v", results)
		}
	})

	t.Run("count query", func(t *testing.T) {
		query := "SELECT COUNT(*) FROM berries"
		results, err := QueryBerryData(query)

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		if len(results) != 1 {
			t.Errorf("QueryBerryData() should return exactly one result for COUNT query, got %d", len(results))
		}

		if results[0] == "0" {
			t.Error("QueryBerryData() COUNT should return more than 0 berries")
		}
	})

	t.Run("specific berry data query", func(t *testing.T) {
		// Try to get all berry names and verify structure
		query := "SELECT name FROM berries ORDER BY name LIMIT 10"
		results, err := QueryBerryData(query)

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		if len(results) == 0 {
			t.Error("QueryBerryData() should return berry names")
		}

		for _, result := range results {
			if result == "" {
				t.Error("QueryBerryData() should not return empty berry names")
			}
		}
	})

	t.Run("validate berries table exists", func(t *testing.T) {
		// Use a simpler query that works with the single-column return format
		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='berries'"
		results, err := QueryBerryData(query)

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		if len(results) == 0 || results[0] != "1" {
			t.Error("QueryBerryData() berries table should exist")
		}
	})

	t.Run("SQL injection protection", func(t *testing.T) {
		maliciousInput := "'; DROP TABLE berries; --"
		query := "SELECT name FROM berries WHERE name = ?"

		results, err := QueryBerryData(query, maliciousInput)
		if err != nil {
			t.Fatalf("Query with malicious input should not error: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("expected no results for malicious input, got %v", results)
		}
	})
}

func TestQueryBerryDataWithVariadicArgs(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		args        []interface{}
		expectError bool
	}{
		{
			name:        "no args",
			query:       "SELECT COUNT(*) FROM berries",
			args:        []interface{}{},
			expectError: false,
		},
		{
			name:        "one arg",
			query:       "SELECT name FROM berries WHERE name LIKE ? LIMIT 1",
			args:        []interface{}{"%a%"},
			expectError: false,
		},
		{
			name:        "multiple args",
			query:       "SELECT name FROM berries WHERE name LIKE ? OR name LIKE ? LIMIT 1",
			args:        []interface{}{"%a%", "%e%"},
			expectError: false,
		},
		{
			name:        "mixed types",
			query:       "SELECT name FROM berries LIMIT ?",
			args:        []interface{}{5},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := QueryBerryData(tt.query, tt.args...)

			if (err != nil) != tt.expectError {
				t.Errorf("QueryBerryData() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && results == nil {
				t.Error("QueryBerryData() should not return nil results on success")
			}
		})
	}
}

func TestEmbeddedDBExists(t *testing.T) {
	if len(embeddedDB) == 0 {
		t.Error("embeddedDB should not be empty")
	}

	if len(embeddedDB) < 16 {
		t.Error("embeddedDB appears to be too small to be a valid SQLite database")
	}

	sqliteHeader := "SQLite format 3"
	if !strings.HasPrefix(string(embeddedDB[:len(sqliteHeader)]), sqliteHeader) {
		t.Error("embeddedDB does not appear to be a valid SQLite database")
	}
}

func TestQueryBerryDataTempFileCleanup(t *testing.T) {
	initialTempFiles := countTempFiles()

	query := "SELECT COUNT(*) FROM berries"
	_, err := QueryBerryData(query)

	if err != nil {
		t.Fatalf("QueryBerryData() error = %v", err)
	}

	finalTempFiles := countTempFiles()

	if finalTempFiles > initialTempFiles {
		t.Errorf("Temporary files not cleaned up properly. Before: %d, After: %d", initialTempFiles, finalTempFiles)
	}
}

func countTempFiles() int {
	tempDir := os.TempDir()
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "berries-") && strings.HasSuffix(entry.Name(), ".db") {
			count++
		}
	}
	return count
}
