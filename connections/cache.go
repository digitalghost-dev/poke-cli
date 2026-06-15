package connections

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var cacheWarnOnce sync.Once

func cacheNotice() (string, error) {
	if suppressCacheWarning() {
		return "", nil
	}
	return "⚠ poke-cache not installed — running without local caching.\n" +
		"  Install it for faster repeat lookups: https://docs.poke-cli.com/caching", nil
}

func warnNoCache() {
	cacheWarnOnce.Do(func() {
		if msg, _ := cacheNotice(); msg != "" {
			fmt.Fprintln(os.Stderr, msg)
		}
	})
}

func suppressCacheWarning() bool {
	return os.Getenv("POKE_CLI_NO_CACHE_WARNING") != ""
}

func cachedFetch(url string) ([]byte, error) {
	if flag.Lookup("test.v") != nil {
		return directFetch(url)
	}
	path, err := exec.LookPath("poke-cache")
	if err != nil {
		warnNoCache()
		return directFetch(url)
	}
	out, err := exec.Command(path, "get", url).Output()
	if err != nil {
		return directFetch(url)
	}
	return out, nil
}