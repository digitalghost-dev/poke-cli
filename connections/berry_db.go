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

func BerryListAllNames() ([]string, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "berries-*.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// temp file
	if _, err := tmpFile.Write(embeddedDB); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write embedded database: %w", err)
	}
	tmpFile.Close()

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

	rows, err := db.Query(`
		SELECT 
		    UPPER(SUBSTR(name, 1, 1)) || SUBSTR(name, 2) 
		FROM 
		    berries 
		ORDER BY 
		    name`)
	if err != nil {
		return nil, fmt.Errorf("failed to query berry names: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan berry name: %w", err)
		}
		names = append(names, name)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return names, nil
}
