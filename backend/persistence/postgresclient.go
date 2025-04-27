package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/etboye/calculator/calculation"
	"github.com/etboye/calculator/errorid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pageMaxNumberOfItems = 5

type PostgresClient struct {
	conn *pgxpool.Pool
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
		return HistoryRow{}, errors.New(errorid.INSERT_CALCULATION_ERROR)
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

func (postgresClient PostgresClient) GetSessionHistoryFirstPage(sessionId string) (CalculationsPageObject, error) {
	cursorPositions, err := postgresClient.computeCursorPositions(sessionId, nil)

	if err != nil {
		return CalculationsPageObject{}, errors.New(errorid.HISTORY_CURSOR_ERROR)
	}

	return postgresClient.GetSessionHistoryFromCursor(sessionId, cursorPositions)
}

func (postgresClient PostgresClient) GetSessionHistory(sessionId string, cursor int64) (CalculationsPageObject, error) {
	cursorPositions, err := postgresClient.computeCursorPositions(sessionId, &cursor)

	if err != nil {
		return CalculationsPageObject{}, errors.New(errorid.HISTORY_CURSOR_ERROR)
	}

	return postgresClient.GetSessionHistoryFromCursor(sessionId, cursorPositions)
}

func (postgresClient PostgresClient) GetSessionHistoryFromCursor(sessionId string, cursorPositions cursorPositions) (CalculationsPageObject, error) {
	if cursorPositions.Self == nil {
		// Empty list, send empty list
		return CalculationsPageObject{Items: []HistoryRow{}}, nil
	}

	log.Printf("Fetching session history for sessionId %s with cursor %d", sessionId, cursorPositions.Self)
	rows, err := postgresClient.conn.Query(context.Background(),
		`SELECT calculationId, input, outputnum, outputdenom, outputestimate, error 
		 FROM HISTORY WHERE sessionId=$1 AND calculationId<=$2 ORDER BY calculationId DESC LIMIT $3`,
		sessionId,
		cursorPositions.Self,
		pageMaxNumberOfItems)

	if err != nil {
		log.Printf("GetSessionHistory: Could not SELECT. Error: %s", err.Error())
		return CalculationsPageObject{}, errors.New(errorid.HISTORY_FETCH_ERROR)
	}

	items := make([]HistoryRow, 0, pageMaxNumberOfItems)

	for rows.Next() {
		calculationsResult, err := scanHistoryRow(rows)

		if err != nil {
			log.Printf("GetSessionHistory: Could not scan. Error: %s", err.Error())
			return CalculationsPageObject{}, errors.New(errorid.HISTORY_SCAN_ROW_ERROR)
		}

		items = append(items, calculationsResult)
	}

	return buildCalculationsPageObject(cursorPositions, items), nil
}

func InitPostgresClient() (PostgresClient, error) {
	databaseUrl := getDatabaseUrl()
	conn, err := pgxpool.New(context.Background(), databaseUrl)
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
	return fmt.Sprintf("postgres://%s:%s@%s:5432/calculator", username, password, host)
}
