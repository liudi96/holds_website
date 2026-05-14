package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"portfolio-desk/internal/industrymetrics"
	"sort"
	"strings"
)

const industriesDir = "data/industries"

type IndustryResearch struct {
	ID               string                    `json:"id"`
	Name             string                    `json:"name"`
	Category         string                    `json:"category,omitempty"`
	Status           string                    `json:"status,omitempty"`
	UpdatedAt        string                    `json:"updatedAt,omitempty"`
	MetricsUpdatedAt string                    `json:"metricsUpdatedAt,omitempty"`
	Summary          string                    `json:"summary,omitempty"`
	Discipline       string                    `json:"discipline,omitempty"`
	Keywords         []string                  `json:"keywords,omitempty"`
	LinkedSymbols    []string                  `json:"linkedSymbols,omitempty"`
	KeyQuestions     []string                  `json:"keyQuestions,omitempty"`
	MetricSourceIDs  []string                  `json:"metricSourceIds,omitempty"`
	Metrics          []IndustryMetric          `json:"metrics,omitempty"`
	CompanyAnalyses  []IndustryCompanyAnalysis `json:"companyAnalyses,omitempty"`
	Notes            []IndustryNote            `json:"notes,omitempty"`
}

type IndustryMetric = industrymetrics.Metric

type IndustryMetricPoint = industrymetrics.Point

type IndustryNote struct {
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Summary string `json:"summary,omitempty"`
}

type IndustryCompanyAnalysis struct {
	Symbol      string          `json:"symbol"`
	Name        string          `json:"name,omitempty"`
	Stance      string          `json:"stance,omitempty"`
	Summary     string          `json:"summary,omitempty"`
	Trends      []IndustryTrend `json:"trends,omitempty"`
	Judgments   []string        `json:"judgments,omitempty"`
	Watchpoints []string        `json:"watchpoints,omitempty"`
}

type IndustryTrend struct {
	Label     string `json:"label"`
	Value     string `json:"value,omitempty"`
	Direction string `json:"direction,omitempty"`
	Note      string `json:"note,omitempty"`
}

func loadIndustries() ([]IndustryResearch, error) {
	entries, err := os.ReadDir(industriesDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	industries := make([]IndustryResearch, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".json") {
			continue
		}
		path := filepath.Join(industriesDir, entry.Name())
		body, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var industry IndustryResearch
		if err := json.Unmarshal(body, &industry); err != nil {
			return nil, err
		}
		industry = normalizeIndustryResearch(industry, entry.Name())
		if strings.TrimSpace(industry.ID) == "" || strings.TrimSpace(industry.Name) == "" {
			continue
		}
		industries = append(industries, industry)
	}

	sort.SliceStable(industries, func(i, j int) bool {
		if industries[i].UpdatedAt != industries[j].UpdatedAt {
			return industries[i].UpdatedAt > industries[j].UpdatedAt
		}
		return industries[i].Name < industries[j].Name
	})
	return industries, nil
}

func normalizeIndustryResearch(industry IndustryResearch, filename string) IndustryResearch {
	industry.ID = strings.TrimSpace(industry.ID)
	if industry.ID == "" {
		industry.ID = strings.TrimSuffix(filename, filepath.Ext(filename))
	}
	industry.Name = strings.TrimSpace(industry.Name)
	industry.Category = strings.TrimSpace(industry.Category)
	industry.Status = strings.TrimSpace(industry.Status)
	industry.UpdatedAt = strings.TrimSpace(industry.UpdatedAt)
	industry.MetricsUpdatedAt = strings.TrimSpace(industry.MetricsUpdatedAt)
	industry.Summary = strings.TrimSpace(industry.Summary)
	industry.Discipline = strings.TrimSpace(industry.Discipline)
	industry.LinkedSymbols = normalizeIndustrySymbols(industry.LinkedSymbols)
	industry.Keywords = normalizeStringList(industry.Keywords)
	industry.KeyQuestions = normalizeStringList(industry.KeyQuestions)
	industry.MetricSourceIDs = normalizeStringList(industry.MetricSourceIDs)
	industry.Metrics = industrymetrics.NormalizeMetrics(industry.Metrics)
	industry.CompanyAnalyses = normalizeIndustryCompanyAnalyses(industry.CompanyAnalyses)
	return industry
}

func normalizeIndustrySymbols(symbols []string) []string {
	normalized := make([]string, 0, len(symbols))
	seen := make(map[string]bool, len(symbols))
	for _, symbol := range symbols {
		value := normalizeSymbol(symbol)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		normalized = append(normalized, value)
	}
	return normalized
}

func normalizeStringList(values []string) []string {
	normalized := make([]string, 0, len(values))
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		text := strings.TrimSpace(value)
		if text == "" || seen[text] {
			continue
		}
		seen[text] = true
		normalized = append(normalized, text)
	}
	return normalized
}

func normalizeIndustryCompanyAnalyses(analyses []IndustryCompanyAnalysis) []IndustryCompanyAnalysis {
	normalized := make([]IndustryCompanyAnalysis, 0, len(analyses))
	seen := make(map[string]bool, len(analyses))
	for _, analysis := range analyses {
		analysis.Symbol = normalizeSymbol(analysis.Symbol)
		analysis.Name = strings.TrimSpace(analysis.Name)
		analysis.Stance = strings.TrimSpace(analysis.Stance)
		analysis.Summary = strings.TrimSpace(analysis.Summary)
		analysis.Judgments = normalizeStringList(analysis.Judgments)
		analysis.Watchpoints = normalizeStringList(analysis.Watchpoints)
		if analysis.Symbol == "" || seen[analysis.Symbol] {
			continue
		}
		seen[analysis.Symbol] = true
		normalized = append(normalized, analysis)
	}
	return normalized
}
