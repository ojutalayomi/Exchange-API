package utils

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fogleman/gg"
)

type ImageData struct {
	TotalCountries  int
	TopCountries    []CountryGDP
	LastRefreshedAt string
}

type CountryGDP struct {
	Name         string
	EstimatedGDP float64
}

func GenerateSummaryImage(countries []CountriesResponse) error {

	cacheDir := "cache"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}

	totalCountries := len(countries)

	var countriesWithGDP []CountryGDP
	for _, country := range countries {
		if country.EstimatedGDP != nil {
			countriesWithGDP = append(countriesWithGDP, CountryGDP{
				Name:         country.Name,
				EstimatedGDP: *country.EstimatedGDP,
			})
		}
	}

	sort.Slice(countriesWithGDP, func(i, j int) bool {
		return countriesWithGDP[i].EstimatedGDP > countriesWithGDP[j].EstimatedGDP
	})

	topCountries := countriesWithGDP
	if len(topCountries) > 5 {
		topCountries = topCountries[:5]
	}

	lastRefreshedAt := time.Now().Format("2006-01-02 15:04:05")
	if len(countries) > 0 {
		lastRefreshedAt = countries[0].LastRefreshedAt
	}

	imageData := ImageData{
		TotalCountries:  totalCountries,
		TopCountries:    topCountries,
		LastRefreshedAt: lastRefreshedAt,
	}

	if err := createImage(imageData); err != nil {
		return fmt.Errorf("failed to create image: %v", err)
	}

	return nil
}

func createImage(data ImageData) error {
	const width = 800
	const height = 600

	dc := gg.NewContext(width, height)

	dc.SetColor(color.RGBA{240, 248, 255, 255})
	dc.Clear()

	dc.SetColor(color.RGBA{25, 25, 112, 255})
	if err := dc.LoadFontFace("", 36); err != nil {

		dc.LoadFontFace("", 36)
	}
	dc.DrawStringAnchored("Countries Summary", float64(width)/2, 50, 0.5, 0.5)

	y := 120
	dc.LoadFontFace("", 24)
	dc.SetColor(color.RGBA{0, 0, 139, 255})
	dc.DrawStringAnchored(fmt.Sprintf("Total Countries: %d", data.TotalCountries), float64(width)/2, float64(y), 0.5, 0.5)

	y += 80
	dc.LoadFontFace("", 20)
	dc.SetColor(color.RGBA{0, 0, 0, 255})
	dc.DrawStringAnchored("Top 5 Countries by Estimated GDP:", float64(width)/2, float64(y), 0.5, 0.5)

	y += 50
	for i, country := range data.TopCountries {
		if i >= 5 {
			break
		}
		dc.LoadFontFace("", 16)
		dc.SetColor(color.RGBA{0, 100, 0, 255})

		gdpStr := formatNumber(country.EstimatedGDP)
		text := fmt.Sprintf("%d. %s - $%s", i+1, country.Name, gdpStr)
		dc.DrawStringAnchored(text, float64(width)/2, float64(y), 0.5, 0.5)
		y += 35
	}

	y += 40
	dc.LoadFontFace("", 14)
	dc.SetColor(color.RGBA{128, 128, 128, 255})
	dc.DrawStringAnchored(fmt.Sprintf("Last Refreshed: %s", data.LastRefreshedAt), float64(width)/2, float64(y), 0.5, 0.5)

	imagePath := filepath.Join("cache", "summary.png")
	if err := dc.SavePNG(imagePath); err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	return nil
}

func formatNumber(num float64) string {
	str := fmt.Sprintf("%.0f", num)

	if len(str) > 3 {
		result := ""
		for i, char := range str {
			if i > 0 && (len(str)-i)%3 == 0 {
				result += ","
			}
			result += string(char)
		}
		return result
	}
	return str
}
