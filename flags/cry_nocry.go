//go:build nocry

package flags

import (
	"fmt"

	cmdutils "github.com/digitalghost-dev/poke-cli/cmd/utils"
)

func CryFlag(endpoint, pokemonName string) error {
	return fmt.Errorf("%s", cmdutils.FormatError(
		"--cry is not supported in the Docker image because audio playback "+
			"requires host device access that isn't reliable in containers.\n\n"+
			"Install poke-cli natively to use this flag — see installation instructions:\n"+
			"  https://github.com/digitalghost-dev/poke-cli#installation",
	))
}
