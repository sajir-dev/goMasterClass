package utils

const (
	CAD = "CAD"
	USD = "USD"
	EUR = "EUR"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CAD, USD, EUR:
		return true
	}

	return false
}
