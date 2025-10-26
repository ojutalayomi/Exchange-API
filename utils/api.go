package utils

import (
	"io"
	"net/http"
)

type ExchangeRates struct {
	Result             string             `json:"result"`
	Provider           string             `json:"provider"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	TimeEolUnix        int64              `json:"time_eol_unix"`
	BaseCode           string             `json:"base_code"`
	Rates              map[string]float64 `json:"rates"`
}

type Country struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Region     string `json:"region"`
	Population int    `json:"population"`
	Currencies []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Flag string `json:"flag"`
}

type CountriesResponse struct {
	Name            string   `json:"name"`
	Capital         string   `json:"capital"`
	Region          string   `json:"region"`
	Population      int      `json:"population"`
	CurrencyCode    *string  `json:"currency_code"`
	ExchangeRate    *float64 `json:"exchange_rate"`
	EstimatedGDP    *float64 `json:"estimated_gdp"`
	FlagURL         string   `json:"flag_url"`
	LastRefreshedAt string   `json:"last_refreshed_at"`
}

func FetchAPI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
