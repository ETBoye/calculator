package calculation

type RationalNumber struct {
	Num      string `json:"num"`
	Denom    string `json:"denom"`
	Estimate string `json:"estimate"` // Maybe float
}

type CalculationResult struct {
	Input   *string         `json:"input"`
	Result  *RationalNumber `json:"result"`
	ErrorId *string         `json:"errorId"`
}

type Calculator interface {
	Compute(input string) CalculationResult
}
