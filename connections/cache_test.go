package connections

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
	t.Setenv("POKE_CLI_NO_CACHE_WARNING", "")
	assert.False(t, suppressCacheWarning(), "empty env var should not suppress")

	t.Setenv("POKE_CLI_NO_CACHE_WARNING", "1")
	assert.True(t, suppressCacheWarning(), "set env var should suppress")
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
