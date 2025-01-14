package util

const (
	USD = "USD"
	CAD = "CAD"
	EUR = "EUR"
)

func IsCurrencyValid(currency string) bool {

	switch currency {
	case USD, CAD, EUR:
		return true
	}
	return false
}
