package connections

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

//go:embed db/berries.db
var embeddedDB []byte

func QueryBerryData(query string, args ...interface{}) ([]string, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "berries-*.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		// Close file first
		if closeErr := tmpFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close temp file: %v\n", closeErr)
		}

		// Then remove it
		if removeErr := os.Remove(tmpFile.Name()); removeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove temp file %s: %v\n", tmpFile.Name(), removeErr)
		}
	}()

	// Write to temp file
	if _, err := tmpFile.Write(embeddedDB); err != nil {
		return nil, fmt.Errorf("failed to write embedded database: %w", err)
	}

	// Open the temp database file
	db, err := sql.Open("sqlite", tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("failed to close database connection: %v\n", err)
		}
	}(db)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query berry data: %w", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var result string
		if err := rows.Scan(&result); err != nil {
			return nil, fmt.Errorf("failed to scan berry data: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}
