package api

import (
	"net/http"
	"strings"

	"github.com/etboye/calculator/calculation"
	"github.com/etboye/calculator/persistence"
)

type ComputeRequest struct {
	Input string `json:"input"`
}

type ComputeResponse struct {
	HistoryRow *persistence.HistoryRow `json:"historyRow"`
	Error      *string                 `json:"error"`
}

type ComputationHandler interface {
	Compute(sessionId string, computeRequest ComputeRequest) SimpleHttpResponse[ComputeResponse]
}

var EMPTY_INPUT_ERROR string = "Input was empty"

type StandardComputationHandler struct {
	calculator        calculation.Calculator
	persistenceClient persistence.PersistenceClient
}

func NewStandardComputationHandler(calculator calculation.Calculator,
	persistenceClient persistence.PersistenceClient) StandardComputationHandler {

	return StandardComputationHandler{
		calculator:        calculator,
		persistenceClient: persistenceClient,
	}
}

func (c StandardComputationHandler) Compute(sessionId string, computeRequest ComputeRequest) SimpleHttpResponse[ComputeResponse] {

	sessionIdValidationErr := validateSessionId(sessionId)
	if sessionIdValidationErr != nil {
		errorId := sessionIdValidationErr.Error()
		return SimpleHttpResponse[ComputeResponse]{ // TODO: Test
			Status:   http.StatusBadRequest,
			Response: ComputeResponse{Error: &errorId},
		}
	}

	input := computeRequest.Input
	if len(strings.TrimSpace(input)) == 0 {
		return SimpleHttpResponse[ComputeResponse]{ // TODO: Test
			Status:   http.StatusBadRequest,
			Response: ComputeResponse{Error: &EMPTY_INPUT_ERROR},
		}
	}

	calculationResult := c.calculator.Compute(input)

	historyRow, persistenceError := c.persistenceClient.SaveComputation(sessionId, input, calculationResult)

	if persistenceError != nil {
		errorId := persistenceError.Error()
		return SimpleHttpResponse[ComputeResponse]{ // TODO: Test
			Status:   http.StatusInternalServerError,
			Response: ComputeResponse{Error: &errorId},
		}
	}

	return SimpleHttpResponse[ComputeResponse]{
		Status: http.StatusCreated,
		Response: ComputeResponse{
			HistoryRow: &historyRow,
		},
	}
}
