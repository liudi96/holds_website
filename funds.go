package main

import (
	"errors"
	"strings"
	"time"
)

const (
	assetTypeFund = "fund"
	fundTypeETF   = "etf"
	fundTypeOTC   = "otc"
)

type Fund struct {
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	FundType          string  `json:"fundType,omitempty"`
	Shares            float64 `json:"shares"`
	Cost              float64 `json:"cost"`
	CurrentPrice      float64 `json:"currentPrice,omitempty"`
	PreviousClose     float64 `json:"previousClose,omitempty"`
	CurrentPriceDate  string  `json:"currentPriceDate,omitempty"`
	PreviousCloseDate string  `json:"previousCloseDate,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	Category          string  `json:"category,omitempty"`
	Notes             string  `json:"notes,omitempty"`
	UpdatedAt         string  `json:"updatedAt,omitempty"`
}

func normalizeFundSymbol(symbol string) string {
	normalized := normalizeSymbol(symbol)
	return strings.TrimSuffix(normalized, ".OF")
}

func splitFundInput(input string) (string, string) {
	text := strings.TrimSpace(input)
	if text == "" {
		return "", ""
	}
	parts := strings.Fields(text)
	if len(parts) > 1 && isFundSymbolLike(parts[0]) {
		return normalizeFundSymbol(parts[0]), strings.TrimSpace(strings.TrimPrefix(text, parts[0]))
	}
	if isFundSymbolLike(text) {
		return normalizeFundSymbol(text), ""
	}
	return normalizeFundSymbol(text), text
}

func isFundSymbolLike(symbol string) bool {
	normalized := strings.TrimSuffix(normalizeSymbol(symbol), ".OF")
	if strings.HasSuffix(normalized, ".SH") || strings.HasSuffix(normalized, ".SZ") || strings.HasSuffix(normalized, ".HK") {
		code := normalized[:len(normalized)-3]
		return allDigits(code) && len(code) >= 4
	}
	return allDigits(normalized) && len(normalized) == 6
}

func allDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func normalizeFundType(fundType string, symbol string) string {
	value := strings.ToLower(strings.TrimSpace(fundType))
	if value == fundTypeETF {
		return fundTypeETF
	}
	normalized := normalizeFundSymbol(symbol)
	if strings.HasSuffix(normalized, ".SH") || strings.HasSuffix(normalized, ".SZ") {
		return fundTypeETF
	}
	return fundTypeOTC
}

func normalizeFund(fund Fund) Fund {
	fund.Symbol = normalizeFundSymbol(fund.Symbol)
	fund.Name = strings.TrimSpace(fund.Name)
	if fund.Name == "" {
		fund.Name = fund.Symbol
	}
	fund.FundType = normalizeFundType(fund.FundType, fund.Symbol)
	fund.Currency = strings.ToUpper(strings.TrimSpace(fund.Currency))
	if fund.Currency == "" {
		fund.Currency = inferCurrencyFromSymbol(fund.Symbol)
	}
	fund.Category = strings.TrimSpace(fund.Category)
	fund.Notes = strings.TrimSpace(fund.Notes)
	fund.CurrentPriceDate = strings.TrimSpace(fund.CurrentPriceDate)
	fund.PreviousCloseDate = strings.TrimSpace(fund.PreviousCloseDate)
	fund.UpdatedAt = strings.TrimSpace(fund.UpdatedAt)
	if fund.CurrentPrice > 0 && fund.PreviousClose <= 0 {
		fund.PreviousClose = fund.CurrentPrice
	}
	return fund
}

func normalizeFunds(funds []Fund) []Fund {
	result := make([]Fund, 0, len(funds))
	seen := map[string]int{}
	for _, fund := range funds {
		fund = normalizeFund(fund)
		if fund.Symbol == "" {
			continue
		}
		if idx, ok := seen[fund.Symbol]; ok {
			result[idx] = mergeFund(result[idx], fund)
			continue
		}
		seen[fund.Symbol] = len(result)
		result = append(result, fund)
	}
	return result
}

func mergeFund(existing Fund, patch Fund) Fund {
	if strings.TrimSpace(patch.Name) != "" {
		existing.Name = strings.TrimSpace(patch.Name)
	}
	if strings.TrimSpace(patch.FundType) != "" {
		existing.FundType = normalizeFundType(patch.FundType, existing.Symbol)
	}
	if patch.Shares > 0 || existing.Shares == 0 {
		existing.Shares = patch.Shares
	}
	if patch.Cost > 0 || existing.Cost == 0 {
		existing.Cost = patch.Cost
	}
	if patch.CurrentPrice > 0 {
		existing.CurrentPrice = patch.CurrentPrice
	}
	if patch.PreviousClose > 0 {
		existing.PreviousClose = patch.PreviousClose
	}
	if strings.TrimSpace(patch.CurrentPriceDate) != "" {
		existing.CurrentPriceDate = strings.TrimSpace(patch.CurrentPriceDate)
	}
	if strings.TrimSpace(patch.PreviousCloseDate) != "" {
		existing.PreviousCloseDate = strings.TrimSpace(patch.PreviousCloseDate)
	}
	if strings.TrimSpace(patch.Currency) != "" {
		existing.Currency = strings.ToUpper(strings.TrimSpace(patch.Currency))
	}
	if strings.TrimSpace(patch.Category) != "" {
		existing.Category = strings.TrimSpace(patch.Category)
	}
	if strings.TrimSpace(patch.Notes) != "" {
		existing.Notes = strings.TrimSpace(patch.Notes)
	}
	if strings.TrimSpace(patch.UpdatedAt) != "" {
		existing.UpdatedAt = strings.TrimSpace(patch.UpdatedAt)
	}
	return normalizeFund(existing)
}

func findFundIndex(funds []Fund, symbol string) int {
	normalized := normalizeFundSymbol(symbol)
	for i := range funds {
		if normalizeFundSymbol(funds[i].Symbol) == normalized {
			return i
		}
	}
	return -1
}

func upsertFund(funds []Fund, fund Fund) []Fund {
	fund = normalizeFund(fund)
	if fund.Symbol == "" {
		return funds
	}
	idx := findFundIndex(funds, fund.Symbol)
	if idx == -1 {
		return append(funds, fund)
	}
	funds[idx] = mergeFund(funds[idx], fund)
	return funds
}

func removeFund(funds []Fund, symbol string) []Fund {
	idx := findFundIndex(funds, symbol)
	if idx == -1 {
		return funds
	}
	return append(funds[:idx], funds[idx+1:]...)
}

func (s *Server) findTradeFund(input string) (Fund, bool) {
	text := strings.TrimSpace(input)
	if text == "" {
		return Fund{}, false
	}
	normalized := normalizeFundSymbol(text)
	for _, fund := range s.state.Funds {
		if normalizeFundSymbol(fund.Symbol) == normalized {
			return normalizeFund(fund), true
		}
	}
	for _, fund := range s.state.Funds {
		if tradeNameMatches(fund.Name, text) {
			return normalizeFund(fund), true
		}
	}
	return Fund{}, false
}

func (s *Server) resolveFundTradeInput(trade *Trade) error {
	input := firstNonEmpty(trade.Symbol, trade.Name)
	if trade.Side == "sell" {
		fund, ok := s.findTradeFund(input)
		if !ok || fund.Shares <= 0 {
			return errors.New("fund position not found")
		}
		trade.Symbol = normalizeFundSymbol(fund.Symbol)
		trade.Name = firstNonEmpty(fund.Name, trade.Name)
		if trade.Currency == "" {
			trade.Currency = fund.Currency
		}
		if trade.CurrentPrice <= 0 {
			trade.CurrentPrice = fund.CurrentPrice
		}
		return nil
	}

	if fund, ok := s.findTradeFund(input); ok {
		trade.Symbol = normalizeFundSymbol(fund.Symbol)
		trade.Name = firstNonEmpty(fund.Name, trade.Name)
		if trade.Currency == "" {
			trade.Currency = fund.Currency
		}
		if trade.CurrentPrice <= 0 {
			trade.CurrentPrice = fund.CurrentPrice
		}
		return nil
	}

	symbolInput, nameInput := splitFundInput(firstNonEmpty(trade.Symbol, trade.Name))
	if !isFundSymbolLike(symbolInput) {
		return errors.New("fund symbol is required")
	}
	trade.Symbol = normalizeFundSymbol(symbolInput)
	if trade.Symbol == "" {
		return errors.New("fund symbol is required")
	}
	if trade.Name == "" || strings.EqualFold(strings.TrimSpace(trade.Name), strings.TrimSpace(input)) {
		trade.Name = firstNonEmpty(nameInput, trade.Symbol)
	}
	if trade.Currency == "" {
		trade.Currency = inferCurrencyFromSymbol(trade.Symbol)
	}
	return nil
}

func applyFundTradeToState(state *AppState, trade Trade) {
	trade.Symbol = normalizeFundSymbol(trade.Symbol)
	idx := findFundIndex(state.Funds, trade.Symbol)
	if idx == -1 {
		state.Funds = append(state.Funds, Fund{
			Symbol:       trade.Symbol,
			Name:         firstNonEmpty(trade.Name, trade.Symbol),
			FundType:     normalizeFundType("", trade.Symbol),
			Shares:       0,
			Cost:         trade.Price,
			CurrentPrice: trade.CurrentPrice,
			Currency:     strings.ToUpper(firstNonEmpty(trade.Currency, inferCurrencyFromSymbol(trade.Symbol))),
		})
		idx = len(state.Funds) - 1
	}

	fund := &state.Funds[idx]
	if trade.Side == "buy" {
		totalCost := fund.Shares*fund.Cost + trade.Shares*trade.Price
		fund.Shares += trade.Shares
		if fund.Shares > 0 {
			fund.Cost = totalCost / fund.Shares
		}
	} else {
		fund.Shares -= trade.Shares
		if fund.Shares < 0 {
			fund.Shares = 0
		}
	}
	if strings.TrimSpace(trade.Name) != "" {
		fund.Name = strings.TrimSpace(trade.Name)
	}
	fund.Symbol = trade.Symbol
	fund.Currency = strings.ToUpper(firstNonEmpty(trade.Currency, fund.Currency, inferCurrencyFromSymbol(trade.Symbol)))
	fund.FundType = normalizeFundType(fund.FundType, fund.Symbol)
	fund.CurrentPrice = trade.CurrentPrice
	if fund.PreviousClose <= 0 {
		fund.PreviousClose = trade.CurrentPrice
	}
	if strings.TrimSpace(fund.CurrentPriceDate) == "" {
		today := time.Now().Format("2006-01-02")
		fund.CurrentPriceDate = today
		fund.PreviousCloseDate = today
	}
	fund.UpdatedAt = time.Now().Format("2006-01-02")
	if trade.Side == "sell" && fund.Shares == 0 {
		state.Funds = append(state.Funds[:idx], state.Funds[idx+1:]...)
	}
	state.Trades = append(state.Trades, trade)
	if trade.Side == "buy" {
		state.Cash -= tradeCashValue(state, trade)
	} else {
		state.Cash += tradeCashValue(state, trade)
	}
}

func reverseFundTradeFromState(state *AppState, trade Trade) {
	trade.Symbol = normalizeFundSymbol(trade.Symbol)
	idx := findFundIndex(state.Funds, trade.Symbol)
	side := strings.ToLower(strings.TrimSpace(trade.Side))
	if side == "buy" {
		if idx >= 0 {
			fund := &state.Funds[idx]
			if fund.Shares <= trade.Shares+0.000001 {
				state.Funds = append(state.Funds[:idx], state.Funds[idx+1:]...)
			} else {
				totalCost := fund.Shares*fund.Cost - trade.Shares*trade.Price
				fund.Shares -= trade.Shares
				if fund.Shares > 0 && totalCost > 0 {
					fund.Cost = totalCost / fund.Shares
				}
			}
		}
		state.Cash += tradeCashValue(state, trade)
	} else if side == "sell" {
		if idx == -1 {
			state.Funds = append(state.Funds, Fund{
				Symbol:       trade.Symbol,
				Name:         firstNonEmpty(trade.Name, trade.Symbol),
				FundType:     normalizeFundType("", trade.Symbol),
				Shares:       0,
				Cost:         trade.Price,
				CurrentPrice: trade.CurrentPrice,
				Currency:     strings.ToUpper(firstNonEmpty(trade.Currency, inferCurrencyFromSymbol(trade.Symbol))),
			})
			idx = len(state.Funds) - 1
		}
		fund := &state.Funds[idx]
		fund.Shares += trade.Shares
		if fund.Cost <= 0 {
			fund.Cost = trade.Price
		}
		if fund.CurrentPrice <= 0 {
			fund.CurrentPrice = trade.CurrentPrice
		}
		if strings.TrimSpace(fund.Currency) == "" {
			fund.Currency = strings.ToUpper(firstNonEmpty(trade.Currency, inferCurrencyFromSymbol(trade.Symbol)))
		}
		if strings.TrimSpace(fund.Name) == "" {
			fund.Name = trade.Name
		}
		state.Cash -= tradeCashValue(state, trade)
	}
	state.Funds = normalizeFunds(state.Funds)
}
