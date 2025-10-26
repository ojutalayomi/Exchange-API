package db

import (
	"database/sql"
	"fmt"
	"os"
	"step2/utils"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

var ()

type MySQL struct {
	db *sql.DB
}

func (m *MySQL) Connect() *MySQL {
	mysqlString := os.Getenv("MYSQL_STRING")
	db, err := sql.Open("mysql", mysqlString)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}
	m.db = db

	_, err = m.db.Exec(`
		CREATE TABLE IF NOT EXISTS countries (
			id INT AUTO_INCREMENT PRIMARY KEY UNIQUE,
			name VARCHAR(255) NOT NULL,
			capital VARCHAR(255),
			region VARCHAR(255),
			population BIGINT NOT NULL,
			currency_code VARCHAR(10),
			exchange_rate DOUBLE,
			estimated_gdp DECIMAL(20,1),
			flag_url VARCHAR(512),
			last_refreshed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table created successfully")
	return m
}

func (m *MySQL) InsertCountries(countries []utils.CountriesResponse) error {
	for _, country := range countries {
		_, err := m.db.Exec(`
			INSERT INTO countries (name, capital, region, population, currency_code, exchange_rate, estimated_gdp, flag_url, last_refreshed_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, country.Name, country.Capital, country.Region, country.Population, country.CurrencyCode, country.ExchangeRate, country.EstimatedGDP, country.FlagURL, country.LastRefreshedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MySQL) UpdateCountry(country utils.CountriesResponse) error {
	_, err := m.db.Exec(`
		UPDATE countries SET name = ?, capital = ?, region = ?, population = ?, currency_code = ?, exchange_rate = ?, estimated_gdp = ?, flag_url = ?, last_refreshed_at = ? WHERE name = ?
	`, country.Name, country.Capital, country.Region, country.Population, country.CurrencyCode, country.ExchangeRate, country.EstimatedGDP, country.FlagURL, country.LastRefreshedAt, country.Name)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQL) GetCountries(region, currency, sort *string) ([]utils.CountriesResponse, error) {
	var (
		args  []interface{}
		conds []string
	)
	query := `
		SELECT id, name, capital, region, population, currency_code, exchange_rate, estimated_gdp, flag_url, last_refreshed_at FROM countries
	`
	if region != nil && *region != "" {
		conds = append(conds, "region = ?")
		args = append(args, *region)
	}
	if currency != nil && *currency != "" {
		conds = append(conds, "currency_code = ?")
		args = append(args, *currency)
	}
	if len(conds) > 0 {
		query += " WHERE " + conds[0]
		for i := 1; i < len(conds); i++ {
			query += " AND " + conds[i]
		}
	}

	if sort != nil && *sort == "gdp_desc" {
		query += " ORDER BY estimated_gdp DESC"
	}
	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var countries []utils.CountriesResponse
	for rows.Next() {
		var country utils.CountriesResponse
		var id int
		var currencyCode sql.NullString
		var exchangeRate sql.NullFloat64
		var estimatedGDP sql.NullFloat64

		err := rows.Scan(&id, &country.Name, &country.Capital, &country.Region, &country.Population, &currencyCode, &exchangeRate, &estimatedGDP, &country.FlagURL, &country.LastRefreshedAt)
		if err != nil {
			return nil, err
		}

		// Convert sql.NullString and sql.NullFloat64 to pointers
		if currencyCode.Valid {
			country.CurrencyCode = &currencyCode.String
		}
		if exchangeRate.Valid {
			country.ExchangeRate = &exchangeRate.Float64
		}
		if estimatedGDP.Valid {
			country.EstimatedGDP = &estimatedGDP.Float64
		}

		countries = append(countries, country)
	}
	return countries, nil
}

func (m *MySQL) GetCountry(name string) (utils.CountriesResponse, error) {
	row := m.db.QueryRow(`
		SELECT id, name, capital, region, population, currency_code, exchange_rate, estimated_gdp, flag_url, last_refreshed_at FROM countries WHERE name = ?
	`, name)
	var country utils.CountriesResponse
	var id int
	var currencyCode sql.NullString
	var exchangeRate sql.NullFloat64
	var estimatedGDP sql.NullFloat64

	err := row.Scan(&id, &country.Name, &country.Capital, &country.Region, &country.Population, &currencyCode, &exchangeRate, &estimatedGDP, &country.FlagURL, &country.LastRefreshedAt)
	if err != nil {
		return utils.CountriesResponse{}, err
	}

	// Convert sql.NullString and sql.NullFloat64 to pointers
	if currencyCode.Valid {
		country.CurrencyCode = &currencyCode.String
	}
	if exchangeRate.Valid {
		country.ExchangeRate = &exchangeRate.Float64
	}
	if estimatedGDP.Valid {
		country.EstimatedGDP = &estimatedGDP.Float64
	}

	return country, nil
}

func (m *MySQL) DeleteCountry(name string) error {
	res, err := m.db.Exec(`
		DELETE FROM countries WHERE name = ?
	`, name)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (m *MySQL) GetStats() (int, error) {
	rows, err := m.db.Query(`
		SELECT COUNT(*) FROM countries
	`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalCountries int
	for rows.Next() {
		err := rows.Scan(&totalCountries)
		if err != nil {
			return 0, err
		}
	}
	return totalCountries, nil
}

func (m *MySQL) UpdateCountries(countries []utils.CountriesResponse) error {
	for _, country := range countries {
		_, err := m.db.Exec(`
			UPDATE countries SET name = ?, capital = ?, region = ?, population = ?, currency_code = ?, exchange_rate = ?, estimated_gdp = ?, flag_url = ?, last_refreshed_at = ? WHERE name = ?
		`, country.Name, country.Capital, country.Region, country.Population, country.CurrencyCode, country.ExchangeRate, country.EstimatedGDP, country.FlagURL, country.LastRefreshedAt, country.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
