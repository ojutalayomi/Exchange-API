package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"step2/db"
	"step2/utils"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	port = "8080"
	addr = "0.0.0.0"
)

func SetupRoutes(database *db.MySQL) *gin.Engine {
	exchangeRatesAPIURL := os.Getenv("EXCHANGE_RATES_API_URL")
	countriesAPIURL := os.Getenv("COUNTRIES_API_URL")

	router := gin.Default()
	router.POST("/countries/refresh", func(c *gin.Context) {
		countries, err := utils.FetchAPI(countriesAPIURL)
		if err != nil {
			c.JSON(503, gin.H{"error": "External data source unavailable", "details": "Could not fetch data from countries API"})
			return
		}

		var Countries []utils.Country
		err = json.Unmarshal([]byte(countries), &Countries)
		if err != nil {
			c.JSON(503, gin.H{"error": "External data source unavailable", "details": "Could not fetch data from countries API"})
			return
		}

		var exchangeRateResponse utils.ExchangeRates
		exchangeRate, err := utils.FetchAPI(exchangeRatesAPIURL)
		if err != nil {
			fmt.Printf("Warning: Could not fetch exchange rates: %v\n", err)
		} else {
			err = json.Unmarshal([]byte(exchangeRate), &exchangeRateResponse)
			if err != nil {
				fmt.Printf("Warning: Could not parse exchange rates: %v\n", err)
			}
		}

		// Get ALL existing countries in ONE query
		existingCountries, err := database.GetCountries()
		if err != nil {
			c.JSON(503, gin.H{"error": "Could not fetch existing countries from database"})
			return
		}

		// Create a map for fast lookup
		existingMap := make(map[string]bool)
		for _, existing := range existingCountries {
			existingMap[existing.Name] = true
		}

		var countriesToInsert []utils.CountriesResponse
		var countriesToUpdate []utils.CountriesResponse

		for _, country := range Countries {
			resp := utils.CountriesResponse{
				Name:            country.Name,
				Capital:         country.Capital,
				Region:          country.Region,
				Population:      country.Population,
				FlagURL:         country.Flag,
				LastRefreshedAt: time.Now().Format(time.RFC3339),
			}

			if len(country.Currencies) == 0 {
				resp.CurrencyCode = nil
				resp.ExchangeRate = nil
				resp.EstimatedGDP = nil
			} else {
				resp.CurrencyCode = &country.Currencies[0].Code

				if rate, exists := exchangeRateResponse.Rates[*resp.CurrencyCode]; exists {
					resp.ExchangeRate = &rate
					multiplier := 1000 + rand.Intn(1001) // 1000-2000
					estimatedGDP := float64(country.Population) * float64(multiplier) / rate
					resp.EstimatedGDP = &estimatedGDP
				} else {
					resp.ExchangeRate = nil
					resp.EstimatedGDP = nil
				}
			}

			// Use map lookup instead of database query
			if existingMap[country.Name] {
				countriesToUpdate = append(countriesToUpdate, resp)
			} else {
				countriesToInsert = append(countriesToInsert, resp)
			}
		}

		// Batch insert new countries
		if len(countriesToInsert) > 0 {
			err = database.InsertCountries(countriesToInsert)
			if err != nil {
				c.JSON(503, gin.H{"error": "Could not insert data into database"})
				return
			}
		}

		// Batch update existing countries
		if len(countriesToUpdate) > 0 {
			err = database.UpdateCountries(countriesToUpdate)
			if err != nil {
				c.JSON(503, gin.H{"error": "Could not update data in database"})
				return
			}
		}

		// Generate image asynchronously (optional optimization)
		go func() {
			allCountries, err := database.GetCountries()
			if err != nil {
				fmt.Printf("Warning: Could not fetch countries for image generation: %v\n", err)
				return
			}
			if err := utils.GenerateSummaryImage(allCountries); err != nil {
				fmt.Printf("Warning: Failed to generate summary image: %v\n", err)
			}
		}()

		c.JSON(204, nil)
	})
	router.GET("/countries", func(c *gin.Context) {
		countries, err := database.GetCountries()
		if err != nil {
			c.JSON(404, gin.H{"error": "Country not found"})
			return
		}
		c.JSON(200, countries)
	})
	router.GET("/countries/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(400, gin.H{"error": "Name is required"})
			return
		}
		country, err := database.GetCountry(name)
		if err != nil {
			c.JSON(404, gin.H{"error": "Country not found"})
			return
		}
		c.JSON(200, country)
	})

	router.DELETE("/countries/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(400, gin.H{"error": "Name is required"})
			return
		}
		err := database.DeleteCountry(name)
		if err != nil {
			c.JSON(404, gin.H{"error": "Country not found"})
			return
		}
		c.JSON(200, gin.H{"message": "Country deleted successfully"})
	})

	router.GET("/stats", func(c *gin.Context) {
		stats, err := database.GetStats()
		if err != nil {
			c.JSON(503, gin.H{"error": "Could not get stats from database"})
			return
		}
		c.JSON(200, gin.H{
			"total_countries":   stats,
			"last_refreshed_at": time.Now().Format(time.RFC3339),
		})
	})

	router.GET("/countries/images", func(c *gin.Context) {
		imagePath := "cache/summary.png"

		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "Summary image not found"})
			return
		}

		c.File(imagePath)
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	return router
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	database := db.MySQL{}
	database.Connect()

	router := SetupRoutes(&database)

	fmt.Printf("\n🚀 Step 2 API server starting on port %s\n", port)
	fmt.Println("📝 API Documentation available at: /")
	fmt.Println("🏥 Health check available at: /health")
	fmt.Printf("🔗 Me endpoint: GET /me\n\n")

	if err := router.Run(fmt.Sprintf("%s:%s", addr, port)); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
