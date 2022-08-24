package exceptions

type CreditCodeLimitationError struct{}

func (c CreditCodeLimitationError) Error() string {
	return "Credit code limitation reached"
}
