package utils

import "strings"

var allowedCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"RUB": true,
	// можно расширить
}

func IsValidCurrency(code string) bool {
	_, ok := allowedCurrencies[strings.ToUpper(code)]
	return ok
}
