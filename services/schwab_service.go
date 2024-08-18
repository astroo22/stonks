package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stonks/models"
)

type SchwabService struct {
	APIKey string
}

func NewSchwabService(apiKey string) *SchwabService {
	return &SchwabService{
		APIKey: apiKey,
	}
}

func (s *SchwabService) GetTickerInfo(ticker string) (*models.TickerInfo, error) {
	// Placeholder URL, replace with the actual Schwab API endpoint
	url := fmt.Sprintf("https://api.schwab.com/v1/markets/quotes/%s", ticker)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	tickerInfo := &models.TickerInfo{
		Ticker: ticker,
		Price:  result["price"].(float64), // Adjust based on the actual API response
	}

	return tickerInfo, nil
}

func (s *SchwabService) GetMultipleTickersInfo(tickers []string) ([]models.TickerInfo, error) {
	var results []models.TickerInfo
	for _, ticker := range tickers {
		info, err := s.GetTickerInfo(ticker)
		if err != nil {
			return nil, err
		}
		results = append(results, *info)
	}
	return results, nil
}
