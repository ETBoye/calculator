package persistence

import "github.com/etboye/calculator/calculation"

type HistoryRow struct {
	CalculationId int64                         `json:"calculationId"`
	Calculation   calculation.CalculationResult `json:"calculation"`
}

type PersistenceClient interface {
	SaveComputation(sessionId string, input string, calculationResult calculation.CalculationResult) (HistoryRow, error)
	GetSessionHistory(sessionId string, cursor int64) (CalculationsPageObject, error)
	GetSessionHistoryFirstPage(sessionId string) (CalculationsPageObject, error)
}

type CalculationsPageObject struct { // Page as in pagination object, not web page
	Self  *int64
	First *int64
	Prev  *int64
	Next  *int64
	Last  *int64
	Items []HistoryRow
}
