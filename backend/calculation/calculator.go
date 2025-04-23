package calculation

type CalculationInput struct {
	Input string
}

type CalculationResult struct {
	Result       *string
	ErrorMessage *string
}

type Calculator interface {
	Compute(input CalculationInput) CalculationResult
}
