package card

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func CardName(cardName string) string {
	return cardName
}

func CardPrice(cardName string) string {
	// URL encode the card name (spaces become %20, / becomes %2F, etc.)
	encodedName := url.QueryEscape(cardName)

	apiURL := fmt.Sprintf("https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/card_pricing_view?select=market_price&number_plus_name=eq.%s", encodedName)
	body, err := CallCardData(apiURL)
	if err != nil {
		return "Price: Not available"
	}

	var results []struct {
		MarketPrice float64 `json:"market_price"`
	}
	err = json.Unmarshal(body, &results)
	if err != nil || len(results) == 0 {
		return "Price: Not available"
	}

	return fmt.Sprintf("Price: $%.2f", results[0].MarketPrice)
}
