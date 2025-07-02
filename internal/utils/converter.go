package utils

import (
	"fmt"
	"math"
)

func ConvertCurrency(from, to string, amount int64) (int64, error) {
	rate, err := GetExchangeRate(from, to)
	if err != nil {
		return 0, err
	}
	converted := float64(amount) * rate
	return int64(math.Round(converted)), nil
}

func GetExchangeRate(from, to string) (float64, error) {
	// TODO: получить из внешнего API (например, exchangeratesapi.io)
	// или использовать заглушку:
	if from == to {
		return 1, nil
	}
	if from == "USD" && to == "EUR" {
		return 0.93, nil
	}
	return 0, fmt.Errorf("no exchange rate for %s -> %s", from, to)
}
