package berry

import (
	"fmt"
	"os"
)

// BerryInfo prints the currently selected berry on a single, updating line.
func BerryInfo(berryName string) {
	// \r     -> move cursor to start of the line
	// \x1b[K -> clear from cursor to end of line
	fmt.Fprintf(os.Stderr, "\r\x1b[KBerry: %s", berryName)
}
