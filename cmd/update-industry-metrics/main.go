package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"portfolio-desk/internal/industrymetrics"
	"strings"
	"time"
)

const (
	defaultDataDir      = "data"
	portfolioDataDirEnv = "PORTFOLIO_DATA_DIR"
)

func defaultDataPath(name string) string {
	dir := strings.TrimSpace(os.Getenv(portfolioDataDirEnv))
	if dir == "" {
		dir = defaultDataDir
	}
	return filepath.Join(filepath.Clean(dir), name)
}

func main() {
	metricsPath := flag.String("metrics", defaultDataPath(filepath.Join("runtime", "industry_metrics.json")), "runtime industry metrics JSON file to update")
	dryRun := flag.Bool("dry-run", false, "fetch and print metrics without writing the runtime file")
	flag.Parse()

	book, skipped, err := industrymetrics.FetchBook(&http.Client{Timeout: 5 * time.Second}, time.Now())
	if err != nil && countMetrics(book) == 0 {
		fail(err)
	}
	if *dryRun {
		body, marshalErr := json.MarshalIndent(struct {
			Updated int                             `json:"updated"`
			Skipped []industrymetrics.SkippedSource `json:"skipped,omitempty"`
			Book    industrymetrics.Book            `json:"book"`
		}{
			Updated: countMetrics(book),
			Skipped: skipped,
			Book:    book,
		}, "", "  ")
		if marshalErr != nil {
			fail(marshalErr)
		}
		fmt.Println(string(body))
		return
	}
	if err := industrymetrics.SaveBook(*metricsPath, book); err != nil {
		fail(err)
	}
	fmt.Printf("updated %d industry metrics in %s\n", countMetrics(book), *metricsPath)
	if len(skipped) > 0 {
		fmt.Printf("skipped %d sources\n", len(skipped))
	}
}

func countMetrics(book industrymetrics.Book) int {
	count := 0
	for _, industry := range book.Industries {
		count += len(industry.Metrics)
	}
	return count
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
