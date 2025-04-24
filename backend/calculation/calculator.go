package calculation

type CalculationInput struct {
	Input string
}

type RationalNumber struct {
	Num      string
	Denom    string
	Estimate string // Maybe float
}

type CalculationResult struct {
	Result  *RationalNumber
	ErrorId *string
}

type Calculator interface {
	Compute(input CalculationInput) CalculationResult
}
