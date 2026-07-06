package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func (s *Server) handleUpsertFund(w http.ResponseWriter, r *http.Request) {
	var patch Fund
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid fund payload")
		return
	}

	pathSymbol := strings.TrimPrefix(r.URL.Path, "/api/funds/")
	if strings.TrimSpace(pathSymbol) != "" && strings.TrimSpace(patch.Symbol) == "" {
		patch.Symbol = pathSymbol
	}
	fund := normalizeFund(patch)
	if fund.Symbol == "" {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if fund.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if fund.CurrentPrice > 0 {
		today := time.Now().Format("2006-01-02")
		if fund.CurrentPriceDate == "" {
			fund.CurrentPriceDate = today
		}
		if fund.PreviousClose <= 0 {
			fund.PreviousClose = fund.CurrentPrice
		}
		if fund.PreviousCloseDate == "" {
			fund.PreviousCloseDate = fund.CurrentPriceDate
		}
	}
	if fund.UpdatedAt == "" {
		fund.UpdatedAt = time.Now().Format("2006-01-02")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	nextState, err := cloneAppState(s.state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	nextState.Funds = upsertFund(nextState.Funds, fund)
	if _, err := saveStateWithBackup(nextState); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save state")
		return
	}
	if err := hydrateState(&nextState); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}
	s.state = nextState
	writeJSON(w, http.StatusOK, s.state)
}

func (s *Server) handleDeleteFund(w http.ResponseWriter, r *http.Request) {
	symbol := normalizeFundSymbol(strings.TrimPrefix(r.URL.Path, "/api/funds/"))
	if symbol == "" {
		writeError(w, http.StatusBadRequest, "missing symbol")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	nextState, err := cloneAppState(s.state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	idx := findFundIndex(nextState.Funds, symbol)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "fund not found")
		return
	}
	if nextState.Funds[idx].Shares > 0 {
		writeError(w, http.StatusConflict, "fund still has shares")
		return
	}
	nextState.Funds = removeFund(nextState.Funds, symbol)
	if _, err := saveStateWithBackup(nextState); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save state")
		return
	}
	if err := hydrateState(&nextState); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}
	s.state = nextState
	writeJSON(w, http.StatusOK, s.state)
}
