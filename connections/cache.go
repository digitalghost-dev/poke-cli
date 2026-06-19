package connections

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/digitalghost-dev/poke-cli/styling"
)

var (
	cacheWarnOnce    sync.Once
	cacheShowWarning = true
	cacheBinaryPath  string
	cacheOnWarn      func()
)

func ConfigureCache(showWarning bool, binaryPath string, onWarn func()) {
	cacheShowWarning = showWarning
	cacheBinaryPath = binaryPath
	cacheOnWarn = onWarn
}

func cacheNotice() (string, error) {
	if suppressCacheWarning() {
		return "", nil
	}
	return styling.WarningColor.Render("poke-cache not installed; running without local caching.\n" +
		"  Install it for faster repeat lookups: https://docs.poke-cli.com/caching"), nil
}

func warnNoCache() {
	cacheWarnOnce.Do(func() {
		msg, _ := cacheNotice()
		if msg == "" {
			return
		}
		fmt.Fprintln(os.Stderr, msg)
		if cacheOnWarn != nil {
			cacheOnWarn()
		}
	})
}

func suppressCacheWarning() bool {
	if v, err := strconv.ParseBool(os.Getenv("POKE_CLI_NO_CACHE_WARNING")); err == nil && v {
		return true
	}
	return !cacheShowWarning
}

func cacheBinary() (string, bool) {
	if cacheBinaryPath != "" {
		if info, err := os.Stat(cacheBinaryPath); err == nil && !info.IsDir() {
			return cacheBinaryPath, true
		}
	}
	path, err := exec.LookPath("poke-cache")
	if err != nil {
		return "", false
	}
	return path, true
}

func cachedFetch(url string) ([]byte, error) {
	if flag.Lookup("test.v") != nil {
		return directFetch(url)
	}
	path, ok := cacheBinary()
	if !ok {
		warnNoCache()
		return directFetch(url)
	}
	out, err := exec.Command(path, "get", url).Output()
	if err != nil {
		return directFetch(url)
	}
	return out, nil
}
