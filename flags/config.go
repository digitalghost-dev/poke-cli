// Package config loads and saves the user's poke-cli preferences as a TOML file.

// config storage across operating systems:
// macOS: ~/Library/Application Support/poke-cli/config.toml
// Windows: %AppData%\poke-cli\config.toml
// Linux: ~/.config/poke-cli/config.toml

package flags

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

const SchemaVersion = 1

const (
	ThemeYellow = "yellow"
	ThemeRed    = "red"
	ThemeBlue   = "blue"
)

const (
	ImageAuto  = "auto"
	ImageKitty = "kitty"
	ImageSixel = "sixel"
	ImageNone  = "none"
)

type Config struct {
	Version int     `toml:"version"`
	Display Display `toml:"display"`
	Cache   Cache   `toml:"cache"`
}

type Display struct {
	Theme         string `toml:"theme"`
	ImageProtocol string `toml:"image_protocol"`
}

// Cache controls the poke-cache integration: whether to show the "not
// installed" warning, and an optional explicit binary path for when poke-cache
// is not on PATH.
type Cache struct {
	ShowWarning bool   `toml:"show_warning"`
	Path        string `toml:"path"`
}

func Defaults() Config {
	return Config{
		Version: SchemaVersion,
		Display: Display{
			Theme:         ThemeYellow,
			ImageProtocol: ImageAuto,
		},
		Cache: Cache{
			ShowWarning: true,
		},
	}
}

// Path resolves where the config lives. os.UserConfigDir gives the OS-native
// base and a poke-cli/ dir is added so future state files have a home alongside it.
func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "poke-cli", "config.toml"), nil
}

// Load resolves the real path and delegates to LoadFrom.
func Load() (Config, bool, error) {
	path, err := Path()
	if err != nil {
		return Defaults(), false, err
	}
	return LoadFrom(path)
}

func LoadFrom(path string) (Config, bool, error) {
	cfg := Defaults()
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, true, nil
	}
	if err != nil {
		return cfg, false, err
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Defaults(), false, err
	}
	cfg.normalize()
	return cfg, false, nil
}

// Save resolves the real path and delegates to SaveTo.
func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	return SaveTo(path, cfg)
}

func SaveTo(path string, cfg Config) error {
	cfg.normalize()
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "config-*.toml")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}

func (c *Config) normalize() {
	switch c.Display.Theme {
	case ThemeYellow, ThemeRed, ThemeBlue:
	default:
		c.Display.Theme = ThemeYellow
	}
	switch c.Display.ImageProtocol {
	case ImageAuto, ImageKitty, ImageSixel, ImageNone:
	default:
		c.Display.ImageProtocol = ImageAuto
	}
	if c.Version == 0 {
		c.Version = SchemaVersion
	}
}
