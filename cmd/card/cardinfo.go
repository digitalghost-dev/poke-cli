package card

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/charmbracelet/x/ansi/sixel"
	"golang.org/x/image/draw"
)

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func CardImage(imageURL string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	parsedURL, err := url.Parse(imageURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return "", fmt.Errorf("invalid URL scheme")
	}
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	limitedBody := io.LimitReader(resp.Body, 10*1024*1024)
	img, _, err := image.Decode(limitedBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	resized := resizeImage(img, 500, 675)

	// Build Sixel string to return
	var buf bytes.Buffer
	buf.WriteString("\x1bPq")
	if err := new(sixel.Encoder).Encode(&buf, resized); err != nil {
		return "", fmt.Errorf("failed to encode sixel: %w", err)
	}
	buf.WriteString("\x1b\\")

	return buf.String(), nil
}
