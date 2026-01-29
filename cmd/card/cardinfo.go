package card

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi/sixel"
	"github.com/dolmen-go/kittyimg"
	"golang.org/x/image/draw"
)

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

// supportsKittyGraphics checks if the terminal supports the Kitty graphics protocol
func supportsKittyGraphics() bool {
	// Check Kitty-specific window ID
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return true
	}

	// Check TERM_PROGRAM for known Kitty-compatible terminals
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	switch termProgram {
	case "kitty", "ghostty", "wezterm":
		return true
	}

	// Check TERM variable for kitty or ghostty
	term := strings.ToLower(os.Getenv("TERM"))
	switch {
	case strings.Contains(term, "kitty"):
		return true
	case strings.Contains(term, "ghostty"):
		return true
	}

	return false
}

// supportsSixelGraphics checks if the terminal supports the Sixel graphics protocol
func supportsSixelGraphics() bool {
	session := os.Getenv("WT_SESSION")
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	term := strings.ToLower(os.Getenv("TERM"))

	if session != "" {
		return true
	}

	// Check TERM_PROGRAM for known Sixel-supporting terminals
	switch termProgram {
	case "iterm.app", "wezterm", "konsole", "tabby", "rio":
		return true
	}

	// Check TERM variable for known Sixel-supporting terminals
	switch {
	case term == "foot" || strings.HasPrefix(term, "foot-"):
		return true
	case term == "xterm-sixel" || strings.Contains(term, "sixel"):
		return true
	}

	return false
}

// CardImage downloads and renders an image using Kitty protocol if supported, otherwise Sixel.
func CardImage(imageURL string) (imageData string, protocol string, err error) {
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	parsedURL, err := url.Parse(imageURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return "", "", errors.New("invalid URL scheme")
	}
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	// Read body into memory first to avoid timeout during decode
	limitedBody := io.LimitReader(resp.Body, 10*1024*1024)
	bodyBytes, err := io.ReadAll(limitedBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to read image data: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(bodyBytes))
	if err != nil {
		return "", "", fmt.Errorf("failed to decode image: %w", err)
	}

	resized := resizeImage(img, 500, 675)

	var buf bytes.Buffer

	if supportsKittyGraphics() {
		if err := kittyimg.Fprint(&buf, resized); err != nil {
			return "", "", fmt.Errorf("failed to encode kitty image: %w", err)
		}
		return buf.String(), "Kitty", nil
	}

	// Fall back to Sixel
	if supportsSixelGraphics() {
		buf.WriteString("\x1bPq")
		if err := new(sixel.Encoder).Encode(&buf, resized); err != nil {
			return "", "", fmt.Errorf("failed to encode sixel: %w", err)
		}
		buf.WriteString("\x1b\\")
		return buf.String(), "Sixel", nil
	}

	// Neither protocol is supported
	return "", "", errors.New("your terminal does not support image rendering (Kitty or Sixel graphics protocols required)\n\nTry using: Kitty, Ghostty, WezTerm, iTerm2, or foot")
}
