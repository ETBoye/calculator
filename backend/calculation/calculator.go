package calculation

type RationalNumber struct {
	Num      string
	Denom    string
	Estimate string // Maybe float
}

type CalculationResult struct {
	Input   *string
	Result  *RationalNumber
	ErrorId *string
}

type Calculator interface {
	Compute(input string) CalculationResult
}
