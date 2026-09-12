package util

const (
	CAD = "CAD"
	EUR = "EUR"
	USD = "USD"
)

func IsSupportedCurrency(currency string )bool {
	switch currency {
	case CAD, EUR, USD:
		return true
	}
	return false
}