package berry

import (
	"fmt"
)

// BerryInfo prints information based on currently selected berry.
func BerryInfo(berryName string) string {
	return fmt.Sprintf("Berry: %s", berryName)
}
