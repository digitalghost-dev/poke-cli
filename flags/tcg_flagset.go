package flags

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/digitalghost-dev/poke-cli/styling"
)

type TcgFlags struct {
	FlagSet  *flag.FlagSet
	Web      *bool
	ShortWeb *bool
}

func SetupTcgFlagSet() *TcgFlags {
	tf := &TcgFlags{}
	tf.FlagSet = flag.NewFlagSet("tcgFlags", flag.ExitOnError)

	tf.Web = tf.FlagSet.Bool("web", false, "Opens a Streamlit dashboard of stats in the browser")
	tf.ShortWeb = tf.FlagSet.Bool("w", false, "Opens a Streamlit dashboard of stats in the browser")

	tf.FlagSet.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"poke-cli tcg [flags]\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-w, --web", "Opens a Streamlit dashboard of stats in the browser."),
		)
		fmt.Println(helpMessage)
	}

	return tf
}

func WebFlag(url string) (string, error) {
	var output strings.Builder

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}

	case "darwin":
		cmd = "open"
		args = []string{url}

	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	if _, err := exec.LookPath(cmd); err != nil {
		fmt.Fprintf(&output, "Can't open a browser in this environment. Visit manually:\n%s\n", url)
		return output.String(), nil
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		fmt.Fprintf(&output, "Failed to open browser: %v\nVisit manually:\n%s\n", err, url)
		return output.String(), nil
	}

	fmt.Fprintf(&output, "Opening: %s\n", url)
	return output.String(), nil
}
