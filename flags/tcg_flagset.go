package flags

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"

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

func WebFlag(url string) error {
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

	return exec.Command(cmd, args...).Start()
}
