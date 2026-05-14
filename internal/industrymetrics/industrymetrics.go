package industrymetrics

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	BrentURL              = "https://fred.stlouisfed.org/graph/fredgraph.csv?id=DCOILBRENTEU"
	BrentMonthlyURL       = "https://datahub.io/core/oil-prices/_r/-/data/brent-monthly.csv"
	DairyIndexURL         = "https://xmsyj.moa.gov.cn/jcyj/"
	DairyFallbackURL      = "https://xmsyj.moa.gov.cn/jcyj/202605/t20260509_6484024.htm"
	DairySupplementPath   = "data/industry_sources/dairy_raw_milk_monthly.json"
	BoyarSearchAPI        = "https://m.boyar.cn/api/article/search"
	ZJBHISearchURL        = "https://www.zjbhi.com/zh/search.html"
	RetailFallbackURL     = "https://www.stats.gov.cn/xxgk/sjfb/zxfb2020/202604/t20260416_1963325.html"
	RealEstateFallbackURL = "https://www.stats.gov.cn/sj/zxfb/202604/t20260416_1963327.html"
	NBSReleaseIndexURL    = "https://www.stats.gov.cn/sj/zxfb/"
)

type Book struct {
	UpdatedAt  string              `json:"updatedAt,omitempty"`
	Industries map[string]Industry `json:"industries"`
}

type Industry struct {
	UpdatedAt string   `json:"updatedAt,omitempty"`
	Metrics   []Metric `json:"metrics,omitempty"`
}

type Metric struct {
	Key         string   `json:"key,omitempty"`
	Name        string   `json:"name"`
	Unit        string   `json:"unit,omitempty"`
	LatestValue *float64 `json:"latestValue,omitempty"`
	ValueText   string   `json:"valueText,omitempty"`
	AsOf        string   `json:"asOf,omitempty"`
	Source      string   `json:"source,omitempty"`
	SourceURL   string   `json:"sourceUrl,omitempty"`
	TrendText   string   `json:"trendText,omitempty"`
	Tone        string   `json:"tone,omitempty"`
	UpdatedAt   string   `json:"updatedAt,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Series      []Point  `json:"series,omitempty"`
}

type Point struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type SkippedSource struct {
	IndustryID string `json:"industryId,omitempty"`
	Key        string `json:"key,omitempty"`
	Source     string `json:"source,omitempty"`
	Error      string `json:"error"`
}

type boyarSearchResponse struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
	Data   struct {
		Pages int `json:"pages"`
		List  []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			AddTime string `json:"addtime"`
		} `json:"list"`
	} `json:"data"`
}

type dairySupplement struct {
	Source    string  `json:"source,omitempty"`
	SourceURL string  `json:"sourceUrl,omitempty"`
	Series    []Point `json:"series"`
}

func LoadBook(path string) (Book, error) {
	book := Book{Industries: map[string]Industry{}}
	body, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return book, nil
	}
	if err != nil {
		return book, err
	}
	if err := json.Unmarshal(body, &book); err != nil {
		return book, err
	}
	return NormalizeBook(book), nil
}

func SaveBook(path string, book Book) error {
	book = NormalizeBook(book)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(path, body, 0o644)
}

func NormalizeBook(book Book) Book {
	book.UpdatedAt = strings.TrimSpace(book.UpdatedAt)
	if book.Industries == nil {
		book.Industries = map[string]Industry{}
	}
	normalized := make(map[string]Industry, len(book.Industries))
	for id, industry := range book.Industries {
		key := NormalizeID(id)
		if key == "" {
			continue
		}
		industry.UpdatedAt = strings.TrimSpace(industry.UpdatedAt)
		industry.Metrics = NormalizeMetrics(industry.Metrics)
		if len(industry.Metrics) == 0 {
			continue
		}
		normalized[key] = industry
	}
	book.Industries = normalized
	return book
}

func NormalizeID(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func NormalizeMetrics(metrics []Metric) []Metric {
	normalized := make([]Metric, 0, len(metrics))
	seen := make(map[string]bool, len(metrics))
	for _, metric := range metrics {
		metric.Key = strings.TrimSpace(metric.Key)
		metric.Name = strings.TrimSpace(metric.Name)
		metric.Unit = strings.TrimSpace(metric.Unit)
		metric.ValueText = strings.TrimSpace(metric.ValueText)
		metric.AsOf = strings.TrimSpace(metric.AsOf)
		metric.Source = strings.TrimSpace(metric.Source)
		metric.SourceURL = strings.TrimSpace(metric.SourceURL)
		metric.TrendText = strings.TrimSpace(metric.TrendText)
		metric.Tone = strings.TrimSpace(metric.Tone)
		metric.UpdatedAt = strings.TrimSpace(metric.UpdatedAt)
		metric.Comment = strings.TrimSpace(metric.Comment)
		if metric.Key == "" {
			metric.Key = metric.Name
		}
		if metric.Name == "" || seen[metric.Key] {
			continue
		}
		seen[metric.Key] = true
		metric.Series = normalizeSeries(metric.Series)
		if metric.LatestValue == nil && len(metric.Series) > 0 {
			value := metric.Series[len(metric.Series)-1].Value
			metric.LatestValue = &value
		}
		normalized = append(normalized, metric)
	}
	sort.SliceStable(normalized, func(i, j int) bool {
		return normalized[i].Key < normalized[j].Key
	})
	return normalized
}

func FetchBook(client *http.Client, now time.Time) (Book, []SkippedSource, error) {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	updatedAt := now.Format(time.RFC3339)
	book := Book{
		UpdatedAt:  updatedAt,
		Industries: map[string]Industry{},
	}
	var skipped []SkippedSource
	add := func(result fetchResult) {
		industryID := result.industryID
		metrics := result.metrics
		err := result.err
		source := result.source
		if err != nil {
			skipped = append(skipped, SkippedSource{IndustryID: industryID, Source: source, Error: err.Error()})
		}
		if len(metrics) == 0 {
			return
		}
		industry := book.Industries[industryID]
		industry.UpdatedAt = updatedAt
		industry.Metrics = append(industry.Metrics, metrics...)
		book.Industries[industryID] = industry
	}

	tasks := []fetchTask{
		{industryID: "oil", source: "FRED / EIA Brent", fetch: fetchOilMetrics},
		{industryID: "dairy", source: "农业农村部畜牧兽医局", fetch: fetchDairyMetrics},
		{industryID: "beverage", source: "国家统计局社零数据", fetch: fetchBeverageMetrics},
		{industryID: "property-services", source: "国家统计局房地产数据", fetch: fetchPropertyMetrics},
	}
	results := make(chan fetchResult, len(tasks))
	for _, task := range tasks {
		task := task
		go func() {
			metrics, err := task.fetch(client, now)
			results <- fetchResult{
				industryID: task.industryID,
				source:     task.source,
				metrics:    metrics,
				err:        err,
			}
		}()
	}
	for range tasks {
		add(<-results)
	}

	book = NormalizeBook(book)
	if metricCount(book) == 0 && len(skipped) > 0 {
		return book, skipped, errors.New("no industry metrics fetched")
	}
	return book, skipped, nil
}

type fetchTask struct {
	industryID string
	source     string
	fetch      func(*http.Client, time.Time) ([]Metric, error)
}

type fetchResult struct {
	industryID string
	source     string
	metrics    []Metric
	err        error
}

func metricCount(book Book) int {
	count := 0
	for _, industry := range book.Industries {
		count += len(industry.Metrics)
	}
	return count
}

func fetchOilMetrics(client *http.Client, now time.Time) ([]Metric, error) {
	body, err := fetchBytes(client, BrentURL)
	sourceURL := BrentURL
	source := "FRED / EIA DCOILBRENTEU"
	if err != nil || len(strings.TrimSpace(string(body))) == 0 {
		body, err = fetchBytes(client, BrentMonthlyURL)
		sourceURL = BrentMonthlyURL
		source = "DataHub oil-prices / EIA Brent monthly"
	}
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewReader(body))
	rows, err := reader.ReadAll()
	if err != nil || len(rows) < 2 {
		return nil, firstNonNil(err, errors.New("empty Brent CSV"))
	}
	var points []Point
	for _, row := range rows[1:] {
		if len(row) < 2 {
			continue
		}
		date := strings.TrimSpace(row[0])
		valueText := strings.TrimSpace(row[1])
		if date == "" || valueText == "" || valueText == "." {
			continue
		}
		value, err := strconv.ParseFloat(valueText, 64)
		if err != nil {
			continue
		}
		points = append(points, Point{Date: date, Value: value})
	}
	series := latestN(monthlyAverage(points), 24)
	metric := metricFromSeries(Metric{
		Key:       "brent-spot",
		Name:      "Brent 原油现货",
		Unit:      "美元/桶",
		Source:    source,
		SourceURL: sourceURL,
		Comment:   "油价是上游油气利润、自由现金流和分红能力的核心周期变量。",
		Series:    series,
	}, "美元/桶", now)
	return []Metric{metric}, nil
}

func fetchDairyMetrics(client *http.Client, now time.Time) ([]Metric, error) {
	links, _ := collectLinks(client, DairyIndexURL, func(text string) bool {
		return strings.Contains(text, "畜产品和饲料集贸市场价格")
	}, 1, 5)
	if len(links) == 0 {
		links = []string{DairyFallbackURL}
	}
	var points []Point
	var lastURL string
	for _, link := range links {
		body, err := fetchText(client, link)
		if err != nil {
			continue
		}
		point, ok := parseDairyPrice(body)
		if !ok {
			continue
		}
		points = append(points, point)
		lastURL = link
		if len(points) >= 140 {
			break
		}
	}
	if extraPoints, sourceURL := fetchBoyarDairyPoints(client); len(extraPoints) > 0 {
		points = append(points, extraPoints...)
		lastURL = sourceURL
	}
	if extraPoints, sourceURL := fetchZJBHIDairyPoints(client); len(extraPoints) > 0 {
		points = append(points, extraPoints...)
		lastURL = sourceURL
	}
	if extraPoints, sourceURL := loadDairySupplementPoints(); len(extraPoints) > 0 {
		points = append(points, extraPoints...)
		lastURL = sourceURL
	}
	if len(points) == 0 {
		return nil, errors.New("failed to parse dairy price")
	}
	series := latestN(monthlyAverage(points), 24)
	metric := metricFromSeries(Metric{
		Key:       "raw-milk-price",
		Name:      "生鲜乳均价",
		Unit:      "元/公斤",
		Source:    "农业农村部公开数据（含博亚和讯/中经百汇转载）",
		SourceURL: firstNonEmpty(lastURL, DairyFallbackURL),
		Comment:   "原奶价格影响乳制品龙头毛利率，但价格低位也说明上游供需偏弱。",
		Series:    series,
	}, "元/公斤", now)
	return []Metric{metric}, nil
}

func loadDairySupplementPoints() ([]Point, string) {
	body, err := os.ReadFile(DairySupplementPath)
	if err != nil {
		return nil, ""
	}
	var supplement dairySupplement
	if err := json.Unmarshal(body, &supplement); err != nil {
		return nil, ""
	}
	return normalizeSeries(supplement.Series), strings.TrimSpace(supplement.SourceURL)
}

func fetchBoyarDairyPoints(client *http.Client) ([]Point, string) {
	var points []Point
	var lastURL string
	seenArticles := map[int]bool{}
	fetchedArticles := 0
	for _, keyword := range []string{"生鲜乳均价", "乳制品市场分析"} {
		pages := 1
		for page := 1; page <= pages && page <= 6; page++ {
			searchURL := boyarSearchURL(keyword, page)
			body, err := fetchBytes(client, searchURL)
			if err != nil {
				continue
			}
			var response boyarSearchResponse
			if err := json.Unmarshal(body, &response); err != nil || response.Errno != 0 {
				continue
			}
			if response.Data.Pages > pages {
				pages = response.Data.Pages
			}
			for _, item := range response.Data.List {
				title := cleanText(item.Title)
				if point, ok := parseDairyPriceFromText(title, 0); ok {
					points = append(points, point)
					lastURL = boyarArticleURL(item.ID)
				}
				if item.ID <= 0 || seenArticles[item.ID] || !strings.Contains(title, "乳制品市场分析") {
					continue
				}
				if fetchedArticles >= 6 {
					continue
				}
				seenArticles[item.ID] = true
				fetchedArticles++
				articleURL := boyarArticleURL(item.ID)
				body, err := fetchText(client, articleURL)
				if err != nil {
					continue
				}
				fallbackYear := parseYearFromText(title)
				for _, point := range parseDairyPricesFromText(cleanText(body), fallbackYear) {
					points = append(points, point)
					lastURL = articleURL
				}
			}
		}
	}
	return normalizeSeries(points), lastURL
}

func fetchZJBHIDairyPoints(client *http.Client) ([]Point, string) {
	var points []Point
	var lastURL string
	for _, keyword := range []string{"生鲜乳价格", "生鲜乳收购价"} {
		for page := 1; page <= 3; page++ {
			searchURL := zjbhiSearchURL(keyword, page)
			body, err := fetchText(client, searchURL)
			if err != nil {
				continue
			}
			items := extractZJBHIResultItems(body)
			if len(items) == 0 {
				items = []string{body}
			}
			for _, item := range items {
				for _, point := range parseDairyPricesFromText(cleanText(item), 0) {
					points = append(points, point)
					lastURL = searchURL
				}
			}
		}
	}
	return normalizeSeries(points), lastURL
}

func fetchBeverageMetrics(client *http.Client, now time.Time) ([]Metric, error) {
	links, _ := collectLinks(client, NBSReleaseIndexURL, func(text string) bool {
		return strings.Contains(text, "社会消费品零售总额")
	}, 4, 18)
	if len(links) == 0 {
		links = []string{RetailFallbackURL}
	}
	amountPoints, yoyPoints, lastURL := parseRetailPages(client, links)
	if len(amountPoints) == 0 && len(yoyPoints) == 0 {
		return nil, errors.New("failed to parse beverage retail metrics")
	}
	metrics := make([]Metric, 0, 2)
	if len(amountPoints) > 0 {
		metrics = append(metrics, metricFromSeries(Metric{
			Key:       "beverage-retail-sales",
			Name:      "饮料类零售额",
			Unit:      "亿元",
			Source:    "国家统计局社会消费品零售总额月度数据",
			SourceURL: firstNonEmpty(lastURL, RetailFallbackURL),
			Comment:   "用限额以上单位饮料类月度零售额观察终端需求方向。",
			Series:    latestN(monthlyLast(amountPoints), 24),
		}, "亿元", now))
	}
	if len(yoyPoints) > 0 {
		metrics = append(metrics, metricFromSeries(Metric{
			Key:       "beverage-retail-yoy",
			Name:      "饮料类零售额同比",
			Unit:      "%",
			Source:    "国家统计局社会消费品零售总额月度数据",
			SourceURL: firstNonEmpty(lastURL, RetailFallbackURL),
			Comment:   "同比方向用于判断饮料终端需求是否连续改善。",
			Series:    latestN(monthlyLast(yoyPoints), 24),
		}, "%", now))
	}
	return metrics, nil
}

func fetchPropertyMetrics(client *http.Client, now time.Time) ([]Metric, error) {
	links, _ := collectLinks(client, NBSReleaseIndexURL, func(text string) bool {
		return strings.Contains(text, "房地产市场基本情况")
	}, 4, 18)
	if len(links) == 0 {
		links = []string{RealEstateFallbackURL}
	}
	sales, completions, funding, lastURL := parsePropertyPages(client, links)
	if len(sales) == 0 && len(completions) == 0 && len(funding) == 0 {
		return nil, errors.New("failed to parse property metrics")
	}
	metrics := make([]Metric, 0, 3)
	if len(sales) > 0 {
		metrics = append(metrics, metricFromSeries(Metric{
			Key:       "property-sales-area-yoy",
			Name:      "商品房销售面积同比",
			Unit:      "%",
			Source:    "国家统计局全国房地产市场基本情况",
			SourceURL: firstNonEmpty(lastURL, RealEstateFallbackURL),
			Comment:   "地产销售弱会间接影响物管新增项目和关联方回款压力。",
			Series:    latestN(monthlyLast(sales), 24),
		}, "%", now))
	}
	if len(completions) > 0 {
		metrics = append(metrics, metricFromSeries(Metric{
			Key:       "property-completion-yoy",
			Name:      "房屋竣工面积同比",
			Unit:      "%",
			Source:    "国家统计局全国房地产市场基本情况",
			SourceURL: firstNonEmpty(lastURL, RealEstateFallbackURL),
			Comment:   "竣工数据影响物管项目交付节奏，但项目质量比单纯面积更重要。",
			Series:    latestN(monthlyLast(completions), 24),
		}, "%", now))
	}
	if len(funding) > 0 {
		metrics = append(metrics, metricFromSeries(Metric{
			Key:       "property-funding-yoy",
			Name:      "房企到位资金同比",
			Unit:      "%",
			Source:    "国家统计局全国房地产市场基本情况",
			SourceURL: firstNonEmpty(lastURL, RealEstateFallbackURL),
			Comment:   "开发资金偏紧时，物管公司应收和关联方回款要重点观察。",
			Series:    latestN(monthlyLast(funding), 24),
		}, "%", now))
	}
	return metrics, nil
}

func metricFromSeries(metric Metric, unit string, now time.Time) Metric {
	metric.Series = normalizeSeries(metric.Series)
	metric.UpdatedAt = now.Format(time.RFC3339)
	if len(metric.Series) == 0 {
		metric.ValueText = "待更新"
		metric.Tone = "watch"
		return metric
	}
	latest := metric.Series[len(metric.Series)-1]
	metric.LatestValue = &latest.Value
	metric.AsOf = latest.Date
	metric.ValueText = valueText(latest.Value, unit)
	metric.TrendText, metric.Tone = trendText(metric.Series, unit)
	return metric
}

func valueText(value float64, unit string) string {
	switch unit {
	case "%":
		return fmt.Sprintf("%.1f%%", value)
	case "亿元":
		return fmt.Sprintf("%.0f亿元", value)
	case "美元/桶":
		return fmt.Sprintf("%.1f美元/桶", value)
	case "元/公斤":
		return fmt.Sprintf("%.2f元/公斤", value)
	default:
		return strconv.FormatFloat(value, 'f', 2, 64)
	}
}

func trendText(series []Point, unit string) (string, string) {
	if len(series) < 2 {
		return "近24个月窗口内只有一个有效点，先展示最新值，等待更多历史数据形成趋势。", "watch"
	}
	first := series[0]
	last := series[len(series)-1]
	change := last.Value - first.Value
	tone := "watch"
	if unit == "%" {
		if change >= 1 {
			tone = "strong"
		} else if change <= -1 {
			tone = "risk"
		}
		return fmt.Sprintf("近24个月窗口：从%s到%s，变化%+.1f个百分点。", valueText(first.Value, unit), valueText(last.Value, unit), change), tone
	}
	if first.Value != 0 {
		changeRatio := change / math.Abs(first.Value)
		if changeRatio >= 0.08 {
			tone = "strong"
		} else if changeRatio <= -0.08 {
			tone = "risk"
		}
		return fmt.Sprintf("近24个月窗口：从%s到%s，变化%+.1f%%。", valueText(first.Value, unit), valueText(last.Value, unit), changeRatio*100), tone
	}
	return fmt.Sprintf("近24个月窗口：从%s到%s。", valueText(first.Value, unit), valueText(last.Value, unit)), tone
}

func parseDairyPrice(body string) (Point, bool) {
	text := cleanText(body)
	date := parseDate(text)
	if date == "" {
		date = parseURLDate(body)
	}
	re := regexp.MustCompile(`生鲜乳价格。.*?生鲜乳平均价格\s*([0-9]+(?:\.[0-9]+)?)\s*元/公斤`)
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		re = regexp.MustCompile(`生鲜乳平均价格\s*([0-9]+(?:\.[0-9]+)?)\s*元/公斤`)
		match = re.FindStringSubmatch(text)
	}
	if len(match) < 2 || date == "" {
		return Point{}, false
	}
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return Point{}, false
	}
	return Point{Date: date, Value: value}, true
}

func parseDairyPriceFromText(text string, fallbackYear int) (Point, bool) {
	points := parseDairyPricesFromText(text, fallbackYear)
	if len(points) == 0 {
		return Point{}, false
	}
	return points[0], true
}

func parseDairyPricesFromText(text string, fallbackYear int) []Point {
	text = cleanText(text)
	var points []Point
	points = append(points, parseDairyPricesWithYear(text)...)
	if fallbackYear > 0 {
		points = append(points, parseDairyPricesWithFallbackYear(text, fallbackYear)...)
	}
	return normalizeSeries(points)
}

func parseDairyPricesWithYear(text string) []Point {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(20\d{2})\s*年\s*(\d{1,2})\s*月[^。；]{0,520}?生鲜乳[^。；]{0,260}?(?:平均)?(?:收购价|价格|均价)[^0-9]{0,30}(?:每公斤\s*)?([0-9]+(?:\.[0-9]+)?)\s*元`),
		regexp.MustCompile(`(20\d{2})\s*年\s*(\d{1,2})\s*月[^。；]{0,260}?生鲜乳[^。；]{0,120}?(?:每公斤\s*)?([0-9]+(?:\.[0-9]+)?)\s*元`),
	}
	var points []Point
	for _, pattern := range patterns {
		for _, match := range pattern.FindAllStringSubmatch(text, -1) {
			if len(match) < 4 {
				continue
			}
			year, _ := strconv.Atoi(match[1])
			month, _ := strconv.Atoi(match[2])
			value, err := strconv.ParseFloat(match[3], 64)
			if err != nil || year < 2000 || month < 1 || month > 12 || !validDairyRawMilkValue(value) {
				continue
			}
			points = append(points, Point{Date: fmt.Sprintf("%04d-%02d", year, month), Value: value})
		}
		if len(points) > 0 {
			return normalizeSeries(points)
		}
	}
	return nil
}

func parseDairyPricesWithFallbackYear(text string, fallbackYear int) []Point {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(\d{1,2})\s*月份?[^。；]{0,180}?生鲜乳[^。；]{0,120}?(?:平均)?(?:收购价|价格|均价)?[^0-9]{0,30}(?:每公斤\s*)?([0-9]+(?:\.[0-9]+)?)\s*元`),
		regexp.MustCompile(`生鲜乳[^。；]{0,120}?(\d{1,2})\s*月份?[^。；]{0,80}?每公斤\s*([0-9]+(?:\.[0-9]+)?)\s*元`),
		regexp.MustCompile(`(\d{1,2})\s*月份?[^。；]{0,24}?每公斤\s*([0-9]+(?:\.[0-9]+)?)\s*元`),
	}
	var points []Point
	for _, re := range patterns {
		for _, match := range re.FindAllStringSubmatch(text, -1) {
			if len(match) < 3 {
				continue
			}
			month, _ := strconv.Atoi(match[1])
			value, err := strconv.ParseFloat(match[2], 64)
			if err != nil || fallbackYear < 2000 || month < 1 || month > 12 || !validDairyRawMilkValue(value) {
				continue
			}
			points = append(points, Point{Date: fmt.Sprintf("%04d-%02d", fallbackYear, month), Value: value})
		}
	}
	return normalizeSeries(points)
}

func validDairyRawMilkValue(value float64) bool {
	return value >= 2 && value <= 5
}

func parseYearFromText(text string) int {
	re := regexp.MustCompile(`(20\d{2})\s*年`)
	match := re.FindStringSubmatch(cleanText(text))
	if len(match) < 2 {
		return 0
	}
	year, _ := strconv.Atoi(match[1])
	return year
}

func extractZJBHIResultItems(body string) []string {
	re := regexp.MustCompile(`(?is)<li>\s*<div class="i">.*?</li>`)
	matches := re.FindAllString(body, -1)
	items := make([]string, 0, len(matches))
	for _, match := range matches {
		if strings.Contains(match, "生鲜乳") {
			items = append(items, match)
		}
	}
	return items
}

func parseRetailPages(client *http.Client, links []string) ([]Point, []Point, string) {
	var amounts []Point
	var yoys []Point
	var lastURL string
	for _, link := range links {
		body, err := fetchText(client, link)
		if err != nil {
			continue
		}
		text := cleanText(body)
		date := parseDate(text)
		re := regexp.MustCompile(`饮料类\s*([0-9]+(?:\.[0-9]+)?)\s*(-?[0-9]+(?:\.[0-9]+)?)`)
		match := re.FindStringSubmatch(text)
		if len(match) < 3 || date == "" {
			continue
		}
		amount, amountErr := strconv.ParseFloat(match[1], 64)
		yoy, yoyErr := strconv.ParseFloat(match[2], 64)
		if amountErr == nil {
			amounts = append(amounts, Point{Date: date, Value: amount})
		}
		if yoyErr == nil {
			yoys = append(yoys, Point{Date: date, Value: yoy})
		}
		lastURL = link
	}
	return amounts, yoys, lastURL
}

func parsePropertyPages(client *http.Client, links []string) ([]Point, []Point, []Point, string) {
	var sales []Point
	var completions []Point
	var funding []Point
	var lastURL string
	for _, link := range links {
		body, err := fetchText(client, link)
		if err != nil {
			continue
		}
		text := cleanText(body)
		date := parseDate(text)
		if date == "" {
			continue
		}
		if value, ok := parsePropertyYoy(text, `新建商品房销售面积（万平方米）`); ok {
			sales = append(sales, Point{Date: date, Value: value})
			lastURL = link
		}
		if value, ok := parsePropertyYoy(text, `房屋竣工面积（万平方米）`); ok {
			completions = append(completions, Point{Date: date, Value: value})
			lastURL = link
		}
		if value, ok := parsePropertyYoy(text, `房地产开发企业本年到位资金（亿元）`); ok {
			funding = append(funding, Point{Date: date, Value: value})
			lastURL = link
		}
	}
	return sales, completions, funding, lastURL
}

func parsePropertyYoy(text string, label string) (float64, bool) {
	re := regexp.MustCompile(regexp.QuoteMeta(label) + `\s*[0-9]+(?:\.[0-9]+)?\s*(-?[0-9]+(?:\.[0-9]+)?)`)
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		return 0, false
	}
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func collectLinks(client *http.Client, indexURL string, accept func(string) bool, maxPages int, targetLinks int) ([]string, error) {
	if maxPages <= 0 {
		maxPages = 1
	}
	seen := map[string]bool{}
	var links []string
	base, err := url.Parse(indexURL)
	if err != nil {
		return nil, err
	}
	for page := 0; page < maxPages; page++ {
		pageURL := pagedIndexURL(indexURL, page)
		body, err := fetchText(client, pageURL)
		if err != nil {
			continue
		}
		for _, link := range extractLinks(body, base) {
			if seen[link.URL] || !accept(link.Text) {
				continue
			}
			seen[link.URL] = true
			links = append(links, link.URL)
			if targetLinks > 0 && len(links) >= targetLinks {
				return links, nil
			}
		}
	}
	return links, nil
}

func boyarSearchURL(keyword string, page int) string {
	values := url.Values{}
	values.Set("keyword", keyword)
	if page > 1 {
		values.Set("page", strconv.Itoa(page))
	}
	return BoyarSearchAPI + "?" + values.Encode()
}

func boyarArticleURL(id int) string {
	return fmt.Sprintf("https://m.boyar.cn/article/%d.html", id)
}

func zjbhiSearchURL(keyword string, page int) string {
	values := url.Values{}
	values.Set("p", "data")
	values.Set("q", keyword)
	if page > 1 {
		values.Set("page", strconv.Itoa(page))
	}
	return ZJBHISearchURL + "?" + values.Encode()
}

type pageLink struct {
	URL  string
	Text string
}

func extractLinks(body string, base *url.URL) []pageLink {
	re := regexp.MustCompile(`(?is)<a[^>]+href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	matches := re.FindAllStringSubmatch(body, -1)
	links := make([]pageLink, 0, len(matches))
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		linkURL := strings.TrimSpace(html.UnescapeString(match[1]))
		text := cleanText(match[2])
		if linkURL == "" || text == "" {
			continue
		}
		resolved, err := base.Parse(linkURL)
		if err != nil {
			continue
		}
		links = append(links, pageLink{URL: resolved.String(), Text: text})
	}
	return links
}

func pagedIndexURL(indexURL string, page int) string {
	if page == 0 {
		return indexURL
	}
	if strings.HasSuffix(indexURL, "/") {
		return fmt.Sprintf("%sindex_%d.html", indexURL, page)
	}
	ext := filepath.Ext(indexURL)
	if ext == "" {
		return fmt.Sprintf("%s/index_%d.html", strings.TrimRight(indexURL, "/"), page)
	}
	return strings.TrimSuffix(indexURL, ext) + fmt.Sprintf("_%d%s", page, ext)
}

func fetchText(client *http.Client, requestURL string) (string, error) {
	body, err := fetchBytes(client, requestURL)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func fetchBytes(client *http.Client, requestURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "holds-website/industry-metrics")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET %s: %s", requestURL, resp.Status)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 8<<20))
}

func cleanText(body string) string {
	reScript := regexp.MustCompile(`(?is)<script.*?</script>|<style.*?</style>`)
	text := reScript.ReplaceAllString(body, " ")
	reTag := regexp.MustCompile(`(?is)<[^>]+>`)
	text = reTag.ReplaceAllString(text, " ")
	text = html.UnescapeString(text)
	reSpace := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(reSpace.ReplaceAllString(text, " "))
}

func parseDate(text string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(20\d{2})[-/](\d{1,2})[-/](\d{1,2})`),
		regexp.MustCompile(`(20\d{2})\s*年\s*1[—\-至到]\s*(\d{1,2})\s*月`),
		regexp.MustCompile(`(20\d{2})\s*年\s*(\d{1,2})\s*月`),
		regexp.MustCompile(`采集日为\s*(\d{1,2})\s*月\s*(\d{1,2})\s*日`),
	}
	for _, pattern := range patterns[:3] {
		match := pattern.FindStringSubmatch(text)
		if len(match) >= 3 {
			year, _ := strconv.Atoi(match[1])
			month, _ := strconv.Atoi(match[2])
			if year >= 2000 && month >= 1 && month <= 12 {
				return fmt.Sprintf("%04d-%02d", year, month)
			}
		}
	}
	match := patterns[3].FindStringSubmatch(text)
	if len(match) >= 3 {
		month, _ := strconv.Atoi(match[1])
		day, _ := strconv.Atoi(match[2])
		year := time.Now().Year()
		if month >= 1 && month <= 12 && day >= 1 && day <= 31 {
			return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		}
	}
	return ""
}

func parseURLDate(text string) string {
	re := regexp.MustCompile(`/((20\d{2})(\d{2}))/`)
	match := re.FindStringSubmatch(text)
	if len(match) >= 4 {
		return fmt.Sprintf("%s-%s", match[2], match[3])
	}
	return ""
}

func monthlyAverage(points []Point) []Point {
	grouped := map[string][]float64{}
	for _, point := range normalizeSeries(points) {
		month := point.Date
		if len(month) >= 7 {
			month = month[:7]
		}
		grouped[month] = append(grouped[month], point.Value)
	}
	series := make([]Point, 0, len(grouped))
	for month, values := range grouped {
		sum := 0.0
		for _, value := range values {
			sum += value
		}
		series = append(series, Point{Date: month, Value: sum / float64(len(values))})
	}
	return normalizeSeries(series)
}

func monthlyLast(points []Point) []Point {
	latest := map[string]Point{}
	for _, point := range normalizeSeries(points) {
		month := point.Date
		if len(month) >= 7 {
			month = month[:7]
		}
		latest[month] = Point{Date: month, Value: point.Value}
	}
	series := make([]Point, 0, len(latest))
	for _, point := range latest {
		series = append(series, point)
	}
	return normalizeSeries(series)
}

func latestN(points []Point, n int) []Point {
	points = normalizeSeries(points)
	if n <= 0 || len(points) <= n {
		return points
	}
	return points[len(points)-n:]
}

func normalizeSeries(points []Point) []Point {
	merged := map[string]Point{}
	for _, point := range points {
		date := strings.TrimSpace(point.Date)
		if date == "" || !isFinite(point.Value) {
			continue
		}
		merged[date] = Point{Date: date, Value: point.Value}
	}
	series := make([]Point, 0, len(merged))
	for _, point := range merged {
		series = append(series, point)
	}
	sort.SliceStable(series, func(i, j int) bool {
		return series[i].Date < series[j].Date
	})
	return series
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstNonNil(values ...error) error {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}
