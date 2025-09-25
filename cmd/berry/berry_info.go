package berry

import (
	"fmt"

	"github.com/digitalghost-dev/poke-cli/connections"
)

// BerryName prints information based on currently selected berry.
func BerryName(berryName string) string {
	return fmt.Sprintf("Berry: %s", berryName)
}

func BerryEffect(berryName string) string {
	berryEffect, err := connections.QueryBerryData(`
		SELECT
			effect
		FROM
			berries
		WHERE
			UPPER(SUBSTR(name, 1, 1)) || SUBSTR(name, 2) = ?`,
		berryName,
	)

	if err != nil || len(berryEffect) == 0 || berryEffect[0] == "" {
		return "Effect information not available"
	}

	return berryEffect[0]
}

func BerryInfo(berryName string) string {
	berryInfo, err := connections.QueryBerryData(`
		SELECT
		   'Firmness: ' || firmness || char(10) ||
		   'Smoothness: ' || smoothness || char(10) ||
		   'Growth Time: ' || growth_time || ' hours' || char(10) ||
		   'Max Harvest: ' || max_harvest
		FROM
			berries
		WHERE
		    UPPER(SUBSTR(name, 1, 1)) || SUBSTR(name, 2) = ?`,
		berryName,
	)

	if err != nil || len(berryInfo) == 0 || berryInfo[0] == "" {
		return "Additional information not available"
	}

	return berryInfo[0]
}
