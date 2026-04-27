package berry

import (
	"image"
	"net/http"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/disintegration/imaging"
)

func berryExists(name string) bool {
	results, err := connections.QueryBerryData(`
		SELECT 1 FROM berries
		WHERE UPPER(SUBSTR(name, 1, 1)) || SUBSTR(name, 2) = ?
		LIMIT 1`,
		name,
	)
	return err == nil && len(results) > 0
}

func berryName(berryName string) string {
	return "Berry: " + berryName
}

func berryEffect(berryName string) string {
	berryEffect, err := connections.QueryBerryData(`
		SELECT
			effect
		FROM
			berries
		WHERE
			UPPER(SUBSTR(name, 1, 1)) || SUBSTR(name, 2) = ?`,
		berryName,
	)

	if err != nil || len(berryEffect) == 0 || berryEffect[0] == "" {
		return "Effect information not available"
	}

	return berryEffect[0]
}

func berryInfo(berryName string) string {
	berryInfo, err := connections.QueryBerryData(`
		SELECT
		   'Firmness: ' || firmness || char(10) ||
		   'Smoothness: ' || smoothness || char(10) ||
		   'Growth Time: ' || growth_time || ' hours' || char(10) ||
		   'Max Harvest: ' || max_harvest
		FROM
			berries
		WHERE
		    UPPER(SUBSTR(name, 1, 1)) || SUBSTR(name, 2) = ?`,
		berryName,
	)

	if err != nil || len(berryInfo) == 0 || berryInfo[0] == "" {
		return "Additional information not available"
	}

	return berryInfo[0]
}

func berryImage(berryName string) string {
	berryImage, err := connections.QueryBerryData(`
		SELECT
			sprite_url
		FROM
			berries
		WHERE
			UPPER(SUBSTR(name, 1, 1)) || SUBSTR(name, 2) = ?`,
		berryName,
	)

	if err != nil || len(berryImage) == 0 || berryImage[0] == "" {
		return "Image information not available"
	}

	ToString := func(width int, height int, img image.Image) string {
		img = imaging.Resize(img, width, height, imaging.NearestNeighbor)
		b := img.Bounds()

		imageWidth := b.Max.X
		h := b.Max.Y

		rowCount := (h - 1) / 2
		if h%2 != 0 {
			rowCount++
		}
		estimatedSize := (imageWidth * rowCount * 10) + rowCount

		str := strings.Builder{}
		str.Grow(estimatedSize)

		styleCache := make(map[string]lipgloss.Style)

		for heightCounter := 0; heightCounter < h-1; heightCounter += 2 {
			for x := 0; x < imageWidth; x++ {
				// Get the color of the current and next row's pixels
				c1, _ := styling.MakeColor(img.At(x, heightCounter))
				color1 := lipgloss.Color(c1.Hex())
				c2, _ := styling.MakeColor(img.At(x, heightCounter+1))
				color2 := lipgloss.Color(c2.Hex())

				styleKey := c1.Hex() + "_" + c2.Hex()
				style, exists := styleCache[styleKey]
				if !exists {
					style = lipgloss.NewStyle().Foreground(color1).Background(color2)
					styleCache[styleKey] = style
				}

				str.WriteString(style.Render("▀"))
			}

			str.WriteString("\n")
		}

		return str.String()
	}

	imageResp, err := http.Get(berryImage[0])
	if err != nil {
		return "Error downloading berry image"
	}
	defer imageResp.Body.Close()

	img, err := imaging.Decode(imageResp.Body)
	if err != nil {
		return "Error decoding berry image"
	}

	return ToString(28, 28, img)
}
