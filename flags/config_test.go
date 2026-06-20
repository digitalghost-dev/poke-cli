package flags

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestDefaults(t *testing.T) {
	cfg := Defaults()

	assert.Equal(t, SchemaVersion, cfg.Version)
	assert.Equal(t, ThemeYellow, cfg.Display.Theme)
	assert.Equal(t, ImageAuto, cfg.Display.ImageProtocol)
	assert.True(t, cfg.Cache.ShowWarning)
}

func TestLoadFromMissingIsFirstRun(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.toml")

	cfg, firstRun, err := LoadFrom(path)

	require.NoError(t, err)
	assert.True(t, firstRun)
	assert.Equal(t, Defaults(), cfg)
}

func TestLoadFromPartialKeepsDefaults(t *testing.T) {
	path := writeTempConfig(t, "[display]\ntheme = \"red\"\n")

	cfg, firstRun, err := LoadFrom(path)

	require.NoError(t, err)
	assert.False(t, firstRun)
	assert.Equal(t, ThemeRed, cfg.Display.Theme)
	assert.Equal(t, ImageAuto, cfg.Display.ImageProtocol)
	assert.True(t, cfg.Cache.ShowWarning)
}

func TestLoadFromCorruptFallsBack(t *testing.T) {
	path := writeTempConfig(t, "[display]\ntheme = \"unterminated\n")

	cfg, firstRun, err := LoadFrom(path)

	require.Error(t, err)
	assert.False(t, firstRun)
	assert.Equal(t, Defaults(), cfg)
}

func TestLoadFromClampsUnknownValues(t *testing.T) {
	path := writeTempConfig(t, "[display]\ntheme = \"chartreuse\"\nimage_protocol = \"ascii\"\n")

	cfg, _, err := LoadFrom(path)

	require.NoError(t, err)
	assert.Equal(t, ThemeYellow, cfg.Display.Theme)
	assert.Equal(t, ImageAuto, cfg.Display.ImageProtocol)
}

func TestSaveToLoadFromRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.toml")
	want := Config{
		Version: SchemaVersion,
		Display: Display{Theme: ThemeBlue, ImageProtocol: ImageKitty},
		Cache:   Cache{ShowWarning: false, Path: "/opt/poke-cache"},
	}

	require.NoError(t, SaveTo(path, want))

	got, firstRun, err := LoadFrom(path)
	require.NoError(t, err)
	assert.False(t, firstRun)
	assert.Equal(t, want, got)
}

func TestSaveToCleansUpTempOnFailure(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config.toml")
	require.NoError(t, os.Mkdir(target, 0o755))

	err := SaveTo(target, Defaults())
	require.Error(t, err)

	entries, readErr := os.ReadDir(dir)
	require.NoError(t, readErr)
	for _, e := range entries {
		assert.False(t, strings.HasPrefix(e.Name(), "config-"), "temp file %q should have been cleaned up", e.Name())
	}
}
