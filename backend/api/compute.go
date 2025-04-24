package api

import (
	"net/http"
	"strings"

	"github.com/etboye/calculator/calculation"
)

type ComputeRequest struct {
	Input string `json:"input"`
}

type ComputeResponse struct {
	CalculationResult *calculation.CalculationResult `json:"calculationResult"`
	Error             *string                        `json:"error"`
}

type ComputationHandler interface {
	GetResponse(computeRequest ComputeRequest) SimpleHttpResponse[ComputeResponse]
}

var EMPTY_INPUT_ERROR string = "Input was empty"

type StandardComputationHandler struct {
	calculator calculation.Calculator
}

func NewStandardComputationHandler(calculator calculation.Calculator) StandardComputationHandler {
	return StandardComputationHandler{
		calculator: calculator,
	}
}

func (c StandardComputationHandler) GetResponse(computeRequest ComputeRequest) SimpleHttpResponse[ComputeResponse] {
	input := computeRequest.Input
	if len(strings.TrimSpace(input)) == 0 {
		return SimpleHttpResponse[ComputeResponse]{ // TODO: Test
			Status:   http.StatusBadRequest,
			Response: ComputeResponse{Error: &EMPTY_INPUT_ERROR},
		}
	}

	calculationResult := c.calculator.Compute(calculation.CalculationInput{Input: input})

	return SimpleHttpResponse[ComputeResponse]{
		Status: http.StatusOK,
		Response: ComputeResponse{
			CalculationResult: &calculationResult,
		},
	}
}
