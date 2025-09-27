package connections

import (
	"os"
	"strings"
	"testing"
)

func TestQueryBerryData(t *testing.T) {
	// Test basic query without parameters
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

	// Test query with parameters
	t.Run("query with parameters", func(t *testing.T) {
		// This test assumes there's a berry in the database
		// We'll use a LIKE query to be more flexible
		query := "SELECT name FROM berries WHERE name LIKE ? LIMIT 1"
		results, err := QueryBerryData(query, "%a%") // Find berries containing 'a'

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		// Should return some results since many berry names contain 'a'
		if len(results) == 0 {
			t.Error("QueryBerryData() should return at least one result for berries containing 'a'")
		}
	})

	// Test query with multiple parameters
	t.Run("query with multiple parameters", func(t *testing.T) {
		query := "SELECT name FROM berries WHERE name LIKE ? OR name LIKE ? LIMIT 5"
		results, err := QueryBerryData(query, "%a%", "%e%")

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		// Should return some results
		if len(results) == 0 {
			t.Error("QueryBerryData() should return results for berries containing 'a' or 'e'")
		}
	})

	// Test invalid query
	t.Run("invalid query", func(t *testing.T) {
		query := "SELECT invalid_column FROM non_existent_table"
		_, err := QueryBerryData(query)

		if err == nil {
			t.Error("QueryBerryData() should return an error for invalid query")
		}
	})

	// Test empty query
	t.Run("empty query", func(t *testing.T) {
		query := ""
		_, err := QueryBerryData(query)

		if err == nil {
			t.Error("QueryBerryData() should return an error for empty query")
		}
	})

	// Test query that returns no results
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

	// Test COUNT query
	t.Run("count query", func(t *testing.T) {
		query := "SELECT COUNT(*) FROM berries"
		results, err := QueryBerryData(query)

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		if len(results) != 1 {
			t.Errorf("QueryBerryData() should return exactly one result for COUNT query, got %d", len(results))
		}

		// The count should be a positive number
		if results[0] == "0" {
			t.Error("QueryBerryData() COUNT should return more than 0 berries")
		}
	})

	// Test specific berry data
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

		// Check that results are strings (berry names)
		for _, result := range results {
			if result == "" {
				t.Error("QueryBerryData() should not return empty berry names")
			}
		}
	})

	// Test database schema validation
	t.Run("validate berries table schema", func(t *testing.T) {
		query := "PRAGMA table_info(berries)"
		results, err := QueryBerryData(query)

		if err != nil {
			t.Fatalf("QueryBerryData() error = %v", err)
		}

		if len(results) == 0 {
			t.Error("QueryBerryData() berries table should exist and have columns")
		}
	})

	// Test SQL injection protection
	t.Run("SQL injection protection", func(t *testing.T) {
		// Try a basic SQL injection attempt
		maliciousInput := "'; DROP TABLE berries; --"
		query := "SELECT name FROM berries WHERE name = ?"

		results, err := QueryBerryData(query, maliciousInput)

		// Should not error (query should execute safely)
		if err != nil {
			// This is okay - the malicious input just won't match anything
			t.Logf("Query with malicious input failed safely: %v", err)
		}

		// Should return empty results since no berry has that name
		if len(results) > 0 {
			t.Errorf("QueryBerryData() should return no results for malicious input, got %v", results)
		}

		// Verify the table still exists by running another query
		testQuery := "SELECT COUNT(*) FROM berries"
		testResults, testErr := QueryBerryData(testQuery)

		if testErr != nil {
			t.Fatalf("Table may have been affected by SQL injection: %v", testErr)
		}

		if len(testResults) == 0 {
			t.Error("Berries table appears to be missing after SQL injection test")
		}
	})
}

func TestQueryBerryDataWithVariadicArgs(t *testing.T) {
	// Test the variadic args functionality specifically
	tests := []struct {
		name     string
		query    string
		args     []interface{}
		expectError bool
	}{
		{
			name:     "no args",
			query:    "SELECT COUNT(*) FROM berries",
			args:     []interface{}{},
			expectError: false,
		},
		{
			name:     "one arg",
			query:    "SELECT name FROM berries WHERE name LIKE ? LIMIT 1",
			args:     []interface{}{"%a%"},
			expectError: false,
		},
		{
			name:     "multiple args",
			query:    "SELECT name FROM berries WHERE name LIKE ? OR name LIKE ? LIMIT 1",
			args:     []interface{}{"%a%", "%e%"},
			expectError: false,
		},
		{
			name:     "mixed types",
			query:    "SELECT name FROM berries LIMIT ?",
			args:     []interface{}{5},
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
	// Test that the embedded database exists and has content
	if len(embeddedDB) == 0 {
		t.Error("embeddedDB should not be empty")
	}

	// Check if it looks like a SQLite database (starts with SQLite magic bytes)
	if len(embeddedDB) < 16 {
		t.Error("embeddedDB appears to be too small to be a valid SQLite database")
	}

	// SQLite files start with "SQLite format 3\000"
	sqliteHeader := "SQLite format 3"
	if !strings.HasPrefix(string(embeddedDB[:len(sqliteHeader)]), sqliteHeader) {
		t.Error("embeddedDB does not appear to be a valid SQLite database")
	}
}

func TestQueryBerryDataTempFileCleanup(t *testing.T) {
	// This test verifies that temporary files are cleaned up properly
	initialTempFiles := countTempFiles()

	// Run a query
	query := "SELECT COUNT(*) FROM berries"
	_, err := QueryBerryData(query)

	if err != nil {
		t.Fatalf("QueryBerryData() error = %v", err)
	}

	finalTempFiles := countTempFiles()

	// The number of temp files should not have increased
	if finalTempFiles > initialTempFiles {
		t.Errorf("Temporary files not cleaned up properly. Before: %d, After: %d", initialTempFiles, finalTempFiles)
	}
}

// Helper function to count temporary files (basic implementation)
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