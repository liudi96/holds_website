package main

import (
	"portfolio-desk/internal/industrymetrics"
	"strings"
)

func loadRuntimeIndustryMetricBook() (industrymetrics.Book, error) {
	return industrymetrics.LoadBook(runtimeIndustryMetricsFile)
}

func saveRuntimeIndustryMetricBook(book industrymetrics.Book) error {
	return industrymetrics.SaveBook(runtimeIndustryMetricsFile, book)
}

func mergeRuntimeIndustryMetrics(industries []IndustryResearch) ([]IndustryResearch, error) {
	book, err := loadRuntimeIndustryMetricBook()
	if err != nil {
		return nil, err
	}
	if len(book.Industries) == 0 {
		return industries, nil
	}
	for i := range industries {
		sourceIDs := append([]string{industries[i].ID}, industries[i].MetricSourceIDs...)
		for _, sourceID := range sourceIDs {
			id := industrymetrics.NormalizeID(sourceID)
			record, ok := book.Industries[id]
			if !ok {
				continue
			}
			industries[i].MetricsUpdatedAt = firstNonEmpty(record.UpdatedAt, book.UpdatedAt, industries[i].MetricsUpdatedAt)
			industries[i].Metrics = mergeIndustryMetrics(industries[i].Metrics, record.Metrics)
		}
	}
	return industries, nil
}

func mergeIndustryMetrics(staticMetrics []IndustryMetric, runtimeMetrics []IndustryMetric) []IndustryMetric {
	merged := make(map[string]IndustryMetric, len(staticMetrics)+len(runtimeMetrics))
	for _, metric := range staticMetrics {
		key := industryMetricKey(metric)
		if key == "" {
			continue
		}
		merged[key] = metric
	}
	for _, metric := range runtimeMetrics {
		key := industryMetricKey(metric)
		if key == "" {
			continue
		}
		merged[key] = metric
	}
	list := make([]IndustryMetric, 0, len(merged))
	for _, metric := range merged {
		list = append(list, metric)
	}
	return industrymetrics.NormalizeMetrics(list)
}

func industryMetricKey(metric IndustryMetric) string {
	key := strings.TrimSpace(metric.Key)
	if key == "" {
		key = strings.TrimSpace(metric.Name)
	}
	return key
}
