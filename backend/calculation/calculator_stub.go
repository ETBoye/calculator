package calculation

type StubCalculator struct{}

func (stubCalculator *StubCalculator) Compute(input CalculationInput) CalculationResult {
	result := "1/1"
	return CalculationResult{
		Result: &result,
	}
}
