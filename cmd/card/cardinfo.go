package card

import (
	"bytes"
	"fmt"
	"image"
	"net/http"

	"github.com/charmbracelet/x/ansi/sixel"
	"golang.org/x/image/draw"
)

func CardName(cardName string) string {
	return cardName
}

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func CardImage(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d\n", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v\n", err)
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
