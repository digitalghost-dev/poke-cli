package connections

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stderr
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	fn()
	require.NoError(t, w.Close())

	out, err := io.ReadAll(r)
	require.NoError(t, err)
	return string(out)
}

func TestSuppressCacheWarning(t *testing.T) {
	tests := []struct {
		value    string
		suppress bool
	}{
		{"", false},
		{"0", false},
		{"false", false},
		{"banana", false},
		{"1", true},
		{"true", true},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Setenv("POKE_CLI_NO_CACHE_WARNING", tt.value)
			assert.Equal(t, tt.suppress, suppressCacheWarning())
		})
	}
}

func TestCacheNotice(t *testing.T) {
	t.Run("not suppressed returns the notice", func(t *testing.T) {
		t.Setenv("POKE_CLI_NO_CACHE_WARNING", "")
		msg, err := cacheNotice()
		require.NoError(t, err)
		assert.Contains(t, msg, "poke-cache not installed")
	})

	t.Run("suppressed returns empty", func(t *testing.T) {
		t.Setenv("POKE_CLI_NO_CACHE_WARNING", "1")
		msg, err := cacheNotice()
		require.NoError(t, err)
		assert.Empty(t, msg)
	})
}

func TestCachedFetch_UnderTestDelegatesToDirectFetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	body, err := cachedFetch(ts.URL)
	require.NoError(t, err)
	assert.Equal(t, `{"ok":true}`, string(body))
}

func TestCachedFetch_PropagatesFetchError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	_, err := cachedFetch(ts.URL)
	require.Error(t, err)

	var statusErr HTTPStatusError
	require.ErrorAs(t, err, &statusErr)
	assert.Equal(t, http.StatusNotFound, statusErr.StatusCode)
}

func TestWarnNoCache_OncePerProcess(t *testing.T) {
	t.Setenv("POKE_CLI_NO_CACHE_WARNING", "")
	cacheWarnOnce = sync.Once{}

	first := captureStderr(t, warnNoCache)
	assert.Contains(t, first, "poke-cache not installed", "expected the notice on the first call")

	second := captureStderr(t, warnNoCache)
	assert.Empty(t, second, "expected nothing on the second call (sync.Once)")
}

func TestConfigureCacheConfigSuppresses(t *testing.T) {
	t.Setenv("POKE_CLI_NO_CACHE_WARNING", "")
	t.Cleanup(func() { ConfigureCache(true, "") })

	ConfigureCache(false, "")
	assert.True(t, suppressCacheWarning(), "show_warning=false should suppress")

	ConfigureCache(true, "")
	assert.False(t, suppressCacheWarning(), "show_warning=true should not suppress")
}

func TestConfigureCacheEnvOverridesConfig(t *testing.T) {
	t.Cleanup(func() { ConfigureCache(true, "") })

	ConfigureCache(true, "")
	t.Setenv("POKE_CLI_NO_CACHE_WARNING", "1")
	assert.True(t, suppressCacheWarning(), "env var should override config")
}

func TestCacheBinaryPrefersConfiguredPath(t *testing.T) {
	t.Cleanup(func() { ConfigureCache(true, "") })

	bin := filepath.Join(t.TempDir(), "poke-cache")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	ConfigureCache(true, bin)
	got, ok := cacheBinary()
	require.True(t, ok)
	assert.Equal(t, bin, got)
}
