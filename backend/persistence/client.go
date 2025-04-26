package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/etboye/calculator/calculation"
	"github.com/jackc/pgx/v5"
)

var pageMaxNumberOfItems = 5

var fetchHistoryErrorId = "HISTORY_FETCH"
var scanRowHistoryError = "HISTORY_SCAN_ROW"
var computeCursorPositionsError = "HISTORY_CURSOR"
var insertCalculationError = "INSERT_CALCULATION"

type HistoryRow struct {
	CalculationId int64                         `json:"calculationId"`
	Calculation   calculation.CalculationResult `json:"calculation"`
}

type CalculationsPageObject struct { // Page as in pagination object, not web page
	Self  *int64
	First *int64
	Prev  *int64
	Next  *int64
	Last  *int64
	Items []HistoryRow
}

type CursorPositions struct {
	Self  *int64
	First *int64
	Prev  *int64
	Next  *int64
	Last  *int64
}

type PersistenceClient interface {
	SaveComputation(sessionId string, input string, calculationResult calculation.CalculationResult) (HistoryRow, error)
	GetSessionHistory(sessionId string, cursor int64) (CalculationsPageObject, error)
	GetSessionHistoryFirstPage(sessionId string) (CalculationsPageObject, error)
}

type PostgresClient struct {
	conn *pgx.Conn
}

func (postgresClient PostgresClient) SaveComputation(
	sessionId string, input string, calculationResult calculation.CalculationResult) (HistoryRow, error) {

	var num, denom, estimate *string
	if calculationResult.Result != nil {
		num = &calculationResult.Result.Num
		denom = &calculationResult.Result.Denom
		estimate = &calculationResult.Result.Estimate
	}
	row := postgresClient.conn.QueryRow(context.Background(),
		`INSERT INTO history (sessionId, input, outputNum, outputDenom, outputEstimate, error) VALUES 
		($1, $2, $3, $4, $5, $6) returning calculationId`,
		sessionId,
		input,
		num, denom, estimate,
		calculationResult.ErrorId,
	)

	var createdCalculationId int64
	err := row.Scan(&createdCalculationId)

	if err != nil {
		return HistoryRow{}, errors.New(insertCalculationError)
	}

	return HistoryRow{
		CalculationId: createdCalculationId,
		Calculation:   calculationResult,
	}, nil
}

func scanHistoryRow(rows pgx.Rows) (HistoryRow, error) {
	var calculationId int64
	var input string
	var outputnum *string
	var outputdenom *string
	var outputestimate *string
	var error *string

	err := rows.Scan(&calculationId, &input, &outputnum, &outputdenom, &outputestimate, &error)

	if err != nil {
		return HistoryRow{}, err
	}

	var rationalNumber *calculation.RationalNumber = nil

	if outputnum != nil && outputdenom != nil && outputestimate != nil {
		rat := calculation.RationalNumber{
			Num:      *outputnum,
			Denom:    *outputdenom,
			Estimate: *outputestimate,
		}

		rationalNumber = &rat
	}

	return HistoryRow{
		CalculationId: calculationId,
		Calculation: calculation.CalculationResult{
			Input:   &input,
			Result:  rationalNumber,
			ErrorId: error,
		},
	}, err
}

func buildCalculationsPageObject(cursorPositions CursorPositions, items []HistoryRow) CalculationsPageObject {
	return CalculationsPageObject{
		Self:  cursorPositions.Self,
		First: cursorPositions.First,
		Prev:  cursorPositions.Prev,
		Next:  cursorPositions.Next,
		Last:  cursorPositions.Last,
		Items: items,
	}
}

func (postgresClient PostgresClient) computeCursorPositions(sessionId string, currentCursorOrNil *int64) (CursorPositions, error) {
	calculationIds, err := postgresClient.getAllCalculationIdsInSessionInDescendingOrder(sessionId)

	if err != nil {
		return CursorPositions{}, err
	}

	if len(calculationIds) == 0 {
		return CursorPositions{
			First: nil,
			Prev:  nil,
			Next:  nil,
			Last:  nil,
		}, nil
	}

	first := calculationIds[0]

	var idxOfLast int
	if len(calculationIds)%pageMaxNumberOfItems == 0 {
		idxOfLast = len(calculationIds) - pageMaxNumberOfItems
	} else {
		idxOfLast = pageMaxNumberOfItems * (len(calculationIds) / pageMaxNumberOfItems)
	}

	last := calculationIds[idxOfLast]

	idxOfCurrent := 0

	if currentCursorOrNil != nil {
		currentCursor := *currentCursorOrNil
		log.Println("computeCursorPositions: Got current cursor", currentCursor)

		for idxOfCurrent < len(calculationIds) && calculationIds[idxOfCurrent] > currentCursor {
			idxOfCurrent++
		}

		log.Println("computeCursorPositions: Computed idxOfCurrent", idxOfCurrent)

		if idxOfCurrent == len(calculationIds) {
			// request is asking for a list with calculation ids below or equal to the lowest calculation id in the list
			// The result will always be an empty list

			log.Printf("Got into special case where idxOfCurrent=%v is out of bounds", idxOfCurrent)

			result := CursorPositions{
				Self:  &currentCursor,
				First: &first,
				Prev:  &last,
				Next:  nil,
				Last:  &last,
			}
			return result, nil
		}

		log.Printf("computeCursorPositions: currentCursor: %d, idxOfCurrent: %d", currentCursor, idxOfCurrent)

	}

	idxOfPrev := idxOfCurrent - pageMaxNumberOfItems
	if idxOfPrev < 0 {
		idxOfPrev = 0
	}

	idxOfNext := idxOfCurrent + pageMaxNumberOfItems

	if idxOfNext >= len(calculationIds) {
		idxOfNext = len(calculationIds) - 1
	}

	self := calculationIds[idxOfCurrent]

	prevResult := &calculationIds[idxOfPrev]
	if self == first {
		prevResult = nil
	}

	nextResult := &calculationIds[idxOfNext]
	if self == last {
		nextResult = nil
	}

	return CursorPositions{
		Self:  &self,
		First: &first,
		Prev:  prevResult,
		Next:  nextResult,
		Last:  &last,
	}, nil
}

func (postgresClient PostgresClient) getAllCalculationIdsInSessionInDescendingOrder(sessionId string) ([]int64, error) {
	rows, err := postgresClient.conn.Query(context.Background(), "SELECT calculationId FROM history WHERE sessionId=$1 ORDER BY calculationId DESC", sessionId)

	if err != nil {
		return []int64{}, err
	}

	result := make([]int64, 0)

	var calculationId int64
	for rows.Next() {
		err := rows.Scan(&calculationId)

		if err != nil {
			return []int64{}, err
		}

		result = append(result, calculationId)
	}

	return result, nil
}

func (postgresClient PostgresClient) GetSessionHistoryFirstPage(sessionId string) (CalculationsPageObject, error) {
	cursorPositions, err := postgresClient.computeCursorPositions(sessionId, nil)

	if err != nil {
		return CalculationsPageObject{}, errors.New(computeCursorPositionsError)
	}

	if cursorPositions.Self == nil {
		// Empty list, send empty list
		return CalculationsPageObject{Items: []HistoryRow{}}, nil
	}
	return postgresClient.GetSessionHistoryFromCursor(sessionId, cursorPositions)
}

func (postgresClient PostgresClient) GetSessionHistory(sessionId string, cursor int64) (CalculationsPageObject, error) {
	cursorPositions, err := postgresClient.computeCursorPositions(sessionId, &cursor)

	if err != nil {
		return CalculationsPageObject{}, errors.New(computeCursorPositionsError)
	}

	if cursorPositions.Self == nil {
		// Empty list, send empty list
		return CalculationsPageObject{Items: []HistoryRow{}}, nil
	}
	return postgresClient.GetSessionHistoryFromCursor(sessionId, cursorPositions)
}

func (postgresClient PostgresClient) GetSessionHistoryFromCursor(sessionId string, cursorPositions CursorPositions) (CalculationsPageObject, error) {
	log.Printf("Fetching session history for sessionId %s with cursor %d", sessionId, cursorPositions.Self)
	rows, err := postgresClient.conn.Query(context.Background(),
		`SELECT calculationId, input, outputnum, outputdenom, outputestimate, error 
		 FROM HISTORY WHERE sessionId=$1 AND calculationId<=$2 ORDER BY calculationId DESC LIMIT $3`,
		sessionId,
		cursorPositions.Self,
		pageMaxNumberOfItems)

	if err != nil {
		log.Printf("GetSessionHistory: Could not SELECT. Error: %s", err.Error())
		return CalculationsPageObject{}, errors.New(fetchHistoryErrorId)
	}

	items := make([]HistoryRow, 0, pageMaxNumberOfItems)

	for rows.Next() {
		calculationsResult, err := scanHistoryRow(rows)

		if err != nil {
			log.Printf("GetSessionHistory: Could not scan. Error: %s", err.Error())
			return CalculationsPageObject{}, errors.New(scanRowHistoryError)
		}

		items = append(items, calculationsResult)
	}

	return buildCalculationsPageObject(cursorPositions, items), nil
}

func InitPostgresClient() (PostgresClient, error) {
	databaseUrl := getDatabaseUrl()
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		return PostgresClient{}, err
	}

	return PostgresClient{conn: conn}, nil
}

func getDatabaseUrl() string {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:5432/over-engineered-calculator", username, password, host)
}
