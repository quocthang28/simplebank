package util

const (
	usd = "USD"
	eur = "EUR"
	cad = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case usd, eur, cad:
		return true
	}

	return false
}
