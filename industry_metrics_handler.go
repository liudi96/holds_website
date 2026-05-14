package main

import (
	"net/http"
	"portfolio-desk/internal/industrymetrics"
	"time"
)

type IndustryMetricsUpdateResponse struct {
	Updated int                             `json:"updated"`
	Skipped []industrymetrics.SkippedSource `json:"skipped,omitempty"`
	State   AppState                        `json:"state"`
}

func (s *Server) handleUpdateIndustries(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, skipped, err := industrymetrics.FetchBook(&http.Client{Timeout: 5 * time.Second}, time.Now())
	if err != nil && countRuntimeIndustryMetrics(book) == 0 {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	if err := saveRuntimeIndustryMetricBook(book); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save runtime industry metrics")
		return
	}
	state, err := loadState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}
	s.state = state
	writeJSON(w, http.StatusOK, IndustryMetricsUpdateResponse{
		Updated: countRuntimeIndustryMetrics(book),
		Skipped: skipped,
		State:   state,
	})
}

func countRuntimeIndustryMetrics(book industrymetrics.Book) int {
	count := 0
	for _, industry := range book.Industries {
		count += len(industry.Metrics)
	}
	return count
}
