package models

import (
	"os"

	"github.com/joho/godotenv"
)

type Client struct {
	apiKey    string
	apiSecret string
	baseURL   string
}
type TickerRequest struct {
	Tickers []string `json:"tickers"`
}

type TickerInfo struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
}

type Settings struct {
	RefreshInterval int      `json:"refresh_interval"`
	Tickers         []string `json:"tickers"`
}

func NewClient() *Client {
	godotenv.Load() // Load .env file

	return &Client{
		apiKey:    os.Getenv("SCHWAB_API_KEY"),
		apiSecret: os.Getenv("SCHWAB_API_SECRET"),
		baseURL:   "https://api.schwab.com/v1", // Adjust base URL as per API documentation
	}
}
