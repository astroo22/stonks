package handlers

import (
	"encoding/json"
	"net/http"
	"stonks/models"
	"stonks/services"
)

type SchwabHandler struct {
	SchwabService *services.SchwabService
	DBService     *services.DBService
}

func NewSchwabHandler(schwabService *services.SchwabService, dbService *services.DBService) *SchwabHandler {
	return &SchwabHandler{
		SchwabService: schwabService,
		DBService:     dbService,
	}
}

func (h *SchwabHandler) GetTickerInfo(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	info, err := h.SchwabService.GetTickerInfo(ticker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.DBService.SaveTickerInfo(*info)
	json.NewEncoder(w).Encode(info)
}

func (h *SchwabHandler) GetMultipleTickersInfo(w http.ResponseWriter, r *http.Request) {
	var req models.TickerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	infos, err := h.SchwabService.GetMultipleTickersInfo(req.Tickers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.DBService.SaveMultipleTickersInfo(infos)
	json.NewEncoder(w).Encode(infos)
}
